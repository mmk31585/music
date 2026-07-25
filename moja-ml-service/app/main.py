from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.api.v1.router import api_router
from app.config import settings

app = FastAPI(
    title=settings.app_name,
    description="Audio/ML processing service for Moja (موجا) music platform",
    version="0.1.0",
)

# ── CORS ───────────────────────────────────────────────────────────────
# This is an internal microservice. Only the Go backend's internal network
# origin should be allowed — never open this to the public internet.
_allowed_origin = (
    settings.go_backend_callback_url.rstrip("/api/v1").rstrip("/")
    if settings.go_backend_callback_url
    else "http://localhost:8080"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=[_allowed_origin],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ── Payload size limit middleware ──────────────────────────────────────
MAX_PAYLOAD_BYTES = 1_048_576  # 1 MB


@app.middleware("http")
async def limit_payload_size(request: Request, call_next):
    if request.method in ("POST", "PUT", "PATCH"):
        content_length = request.headers.get("content-length")
        if content_length and int(content_length) > MAX_PAYLOAD_BYTES:
            return JSONResponse(
                status_code=413,
                content={"detail": "Request body too large (max 1 MB)"},
            )
    return await call_next(request)


# ── Routers ────────────────────────────────────────────────────────────
app.include_router(api_router, prefix="/api/v1")
