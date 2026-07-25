# ML Service Task Refactoring Guide

## Quick Start: How to Refactor a Task Module

### Step 1: Import the Base Classes

```python
from app.workers.tasks.base import BaseTaskOrchestrator, StorageResolutionMixin
from sqlalchemy.orm import Session as SASession
```

### Step 2: Create Your Orchestrator

```python
class YourTaskOrchestrator(BaseTaskOrchestrator[YourJobModel], StorageResolutionMixin):
    """Orchestrator for your specific job type."""

    def _get_job_model(self) -> type[YourJobModel]:
        """Return the SQLAlchemy model class."""
        return YourJobModel

    def _process_job(self, session: SASession, job: YourJobModel) -> None:
        """Implement your business logic here."""
        # 1. Resolve files (if needed)
        local_path = self.resolve_file(
            download_url=job.download_url,
            file_path=job.file_path,
            temp_dir="/tmp/processing",
            settings=settings,
        )
        
        # 2. Update status
        self._set_status(session, job, "processing")
        session.commit()
        
        # 3. Do your work
        result = process_your_data(local_path)
        
        # 4. Store results
        save_results(session, job, result)
        
        # 5. Mark complete (optional - base class does this automatically)
        self._set_status(session, job, "completed")
```

### Step 3: Simplify Your Task Function

```python
@_celery_app.task(
    bind=True,
    max_retries=2,
    default_retry_delay=60,
    name="app.workers.tasks.your_task",
    acks_late=True,
    reject_on_worker_lost=True,
)
def your_task(self: Task, job_id: str) -> None:
    """Your task description."""
    orchestrator = YourTaskOrchestrator(celery_task=self, job_id=job_id)
    orchestrator.execute()
```

### Step 4: Remove Old Helper Functions

Delete these functions (now in base class):
- `_set_status()`
- `_fail_job()`
- `_bump_retry()`
- Manual cleanup in `finally` blocks

## Feature Guide

### Automatic Error Classification

The base orchestrator automatically classifies errors:

```python
# NO NEED TO WRITE THIS - It's automatic!

try:
    self._process_job(session, job)
except CallbackRejectedError:
    # 4xx from Go → permanent failure
    self._fail_job(session, job, str(exc))
except (CallbackTransientError, OSError, ConnectionError):
    # Network issues → retry
    self._bump_retry(session, job, str(exc))
    raise self.task.retry(exc=exc)
except Exception:
    # Processing errors → permanent failure
    self._fail_job(session, job, str(exc))
```

### Cleanup Hooks

Register cleanup functions that run automatically:

```python
def _process_job(self, session: SASession, job: YourJobModel) -> None:
    temp_file = download_file()
    
    # Register cleanup - runs even if processing fails
    self.register_cleanup(lambda: os.remove(temp_file))
    
    # Do work...
```

### Storage Resolution

Use the mixin to handle both storage modes:

```python
# Automatic handling of shared_volume vs presigned_url
local_path = self.resolve_file(
    download_url=job.download_url,  # For presigned_url mode
    file_path=job.file_path,        # For shared_volume mode
    temp_dir="/tmp",
    settings=settings,
)
# Cleanup registered automatically for presigned_url downloads
```

### Optional: Failed Callback

Override to send failure notifications to Go:

```python
class YourTaskOrchestrator(BaseTaskOrchestrator[YourJobModel]):
    # ... other methods ...
    
    def _send_failed_callback(
        self, session: SASession, job: YourJobModel, error_message: str
    ) -> None:
        """Send failure notification to Go backend."""
        client = get_callback_client()
        client.notify_failure(
            job_id=job.id,
            error=error_message,
        )
```

## Real Example: Embedding Tasks

### Before (264 lines)

```python
def extract_embedding_task(self: Task, job_id: str) -> None:
    local_temp_path: str | None = None
    
    with get_sync_db() as session:
        job = session.get(EmbeddingJob, job_id)
        if job is None:
            logger.error("Job not found")
            return
        
        try:
            # Status management
            _set_status(session, job, "downloading")
            session.commit()
            
            # Storage resolution
            storage_client = get_storage_client(settings)
            source = job.audio_file_path or job.audio_download_url
            if not source:
                raise ValueError("No source")
            temp_dir = os.path.join(settings.audio_shared_path, "temp")
            os.makedirs(temp_dir, exist_ok=True)
            local_path = storage_client.resolve(source, dest_dir=temp_dir)
            local_temp_path = local_path
            
            # Processing
            _set_status(session, job, "processing")
            session.commit()
            extractor = get_embedding_extractor(settings)
            service = AudioFeatureExtractionService(extractor=extractor)
            result = service.extract_all(local_path)
            
            # Save results
            _set_status(session, job, "upserting")
            session.commit()
            session.execute(text("INSERT ... ON CONFLICT ..."), {...})
            
            # Complete
            _set_status(session, job, "completed")
            session.commit()
            
        except (AudioEmbeddingError, AudioFeatureExtractionError) as exc:
            logger.warning("Processing failed: %s", exc)
            _fail_job(session, job, str(exc))
        except (OSError, ConnectionError) as exc:
            logger.warning("Retry %d/%d: %s", ...)
            _bump_retry(session, job, str(exc))
            session.commit()
            raise self.retry(exc=exc) from exc
        else:
            if job.status != "completed":
                _set_status(session, job, "completed")
            session.commit()
        finally:
            if local_temp_path and settings.storage_mode == "presigned_url":
                try:
                    os.remove(local_temp_path)
                except OSError:
                    logger.warning("Cleanup failed")

def _set_status(session, job: EmbeddingJob, new_status: str) -> None:
    job.status = new_status
    session.add(job)

def _fail_job(session, job: EmbeddingJob, error_message: str) -> None:
    job.status = "failed"
    job.error_message = error_message
    session.add(job)
    session.commit()

def _bump_retry(session, job: EmbeddingJob, error_message: str) -> None:
    job.retry_count = (job.retry_count or 0) + 1
    job.error_message = error_message
    session.add(job)
```

### After (148 lines)

```python
class EmbeddingTaskOrchestrator(BaseTaskOrchestrator[EmbeddingJob], StorageResolutionMixin):
    """Orchestrator for embedding extraction jobs."""

    def _get_job_model(self) -> type[EmbeddingJob]:
        return EmbeddingJob

    def _process_job(self, session: SASession, job: EmbeddingJob) -> None:
        """Execute embedding extraction pipeline."""
        # Resolve audio
        self._set_status(session, job, "downloading")
        session.commit()
        
        temp_dir = os.path.join(settings.audio_shared_path, "temp-downloads")
        os.makedirs(temp_dir, exist_ok=True)
        local_path = self.resolve_file(
            download_url=job.audio_download_url,
            file_path=job.audio_file_path,
            temp_dir=temp_dir,
            settings=settings,
        )
        
        # Extract features
        self._set_status(session, job, "processing")
        session.commit()
        
        extractor = get_embedding_extractor(settings)
        service = AudioFeatureExtractionService(extractor=extractor)
        result = service.extract_all(local_path)
        
        # Save results
        self._set_status(session, job, "upserting")
        session.commit()
        
        session.execute(text("INSERT ... ON CONFLICT ..."), {...})
        
        self._set_status(session, job, "completed")
        session.commit()

@_celery_app.task(...)
def extract_embedding_task(self: Task, job_id: str) -> None:
    """Extract audio embedding for a track."""
    orchestrator = EmbeddingTaskOrchestrator(celery_task=self, job_id=job_id)
    orchestrator.execute()
```

## Benefits Summary

| Aspect | Before | After |
|--------|--------|-------|
| Lines of code | 264 | 148 |
| Error handling | Manual, duplicated | Automatic, centralized |
| Cleanup | Manual finally block | Cleanup hooks |
| Status transitions | 3 helper functions | Inherited methods |
| Storage resolution | Inline logic | Mixin method |
| Testability | Mock entire task | Mock orchestrator |
| Type safety | Dynamic | Generic types |

## Common Patterns

### Pattern 1: Multi-Step Processing

```python
def _process_job(self, session: SASession, job: JobModel) -> None:
    self._set_status(session, job, "step1")
    session.commit()
    result1 = do_step1()
    
    self._set_status(session, job, "step2")
    session.commit()
    result2 = do_step2(result1)
    
    self._set_status(session, job, "completed")
    session.commit()
```

### Pattern 2: With Callback

```python
def _process_job(self, session: SASession, job: JobModel) -> None:
    result = process()
    
    # Send callback
    client = get_callback_client()
    client.deliver_result(job.id, result)
    
    self._set_status(session, job, "completed")
    session.commit()

def _send_failed_callback(
    self, session: SASession, job: JobModel, error: str
) -> None:
    client = get_callback_client()
    client.deliver_failure(job.id, error)
```

### Pattern 3: Multiple File Processing

```python
def _process_job(self, session: SASession, job: JobModel) -> None:
    # Process multiple files
    for file_url in job.file_urls:
        local_path = self.resolve_file(
            download_url=file_url,
            file_path=None,
            temp_dir="/tmp",
            settings=settings,
        )
        process_file(local_path)
        # Cleanup registered automatically per file
```

## Testing

### Unit Test Structure

```python
def test_process_job():
    # Arrange
    mock_task = Mock()
    orchestrator = YourTaskOrchestrator(
        celery_task=mock_task,
        job_id="test-id"
    )
    mock_session = Mock()
    mock_job = Mock()
    
    # Act
    orchestrator._process_job(mock_session, mock_job)
    
    # Assert
    assert mock_job.status == "completed"
```

## Migration Checklist

- [ ] Create orchestrator class extending `BaseTaskOrchestrator`
- [ ] Add `StorageResolutionMixin` if task handles files
- [ ] Implement `_get_job_model()` method
- [ ] Move business logic to `_process_job()` method
- [ ] Replace task function body with orchestrator instantiation
- [ ] Remove `_set_status`, `_fail_job`, `_bump_retry` functions
- [ ] Remove manual error handling try/except blocks
- [ ] Remove manual cleanup in finally block
- [ ] Update tests to work with new structure
- [ ] Verify task produces same results as before

## FAQ

**Q: Can I customize error handling?**  
A: Override `execute()` but call `super().execute()` for base behavior.

**Q: What if I need custom cleanup logic?**  
A: Use `self.register_cleanup(your_function)` in `_process_job()`.

**Q: How do I handle job-specific status values?**  
A: Use `self._set_status(session, job, "your_custom_status")` freely.

**Q: Can I use this with async tasks?**  
A: Currently sync only. Async version would need separate base class.

**Q: What about tasks that don't follow the standard pattern?**  
A: Keep them as-is or extend the base class minimally.
