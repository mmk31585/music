from fastapi import APIRouter, Depends

from app.api.v1.auth import require_api_key
from app.api.v1.covers import router as cover_router
from app.api.v1.health import router as health_router
from app.api.v1.lyrics import router as lyrics_router
from app.api.v1.similarity import emb_router, sim_router
from app.api.v1.videos import router as video_router

api_router = APIRouter()

api_router.include_router(health_router, tags=["health"])
api_router.include_router(lyrics_router, dependencies=[Depends(require_api_key)])
api_router.include_router(cover_router, dependencies=[Depends(require_api_key)])
api_router.include_router(sim_router, dependencies=[Depends(require_api_key)])
api_router.include_router(emb_router, dependencies=[Depends(require_api_key)])
api_router.include_router(video_router, dependencies=[Depends(require_api_key)])
