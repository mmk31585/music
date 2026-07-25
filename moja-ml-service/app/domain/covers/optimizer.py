"""
Cover image optimizer using Pillow.

Transforms cover art into WebP with proper sizing and EXIF stripping.
Three variants are produced per original image:

+----------------+----------+---------------------+
| Variant        | Max Size | Usage               |
+----------------+----------+---------------------+
| original.webp  | 1200x1200| Full-res fallback   |
| med.webp       | 300x300  | Album cards, lists  |
| thumb.webp     | 100x100  | Mini player, search |
+----------------+----------+---------------------+

The Go backend previously used Go's image/jpeg at quality 85 with
draw.ApproxBiLinear scaling. This module produces smaller files with
better visual quality via Pillow's ``LANCOS`` resampler and WebP's
superior compression.
"""

from __future__ import annotations

import io
import logging
from dataclasses import dataclass

from PIL import Image, ImageOps

logger = logging.getLogger(__name__)

# ── Constants ──────────────────────────────────────────────────────────

# Max dimensions for each variant (preserve aspect ratio)
ORIGINAL_MAX = 1200  # cap at 1200px to avoid serving multi-MB files
MED_MAX = 300        # medium card/list thumbnails
THUMB_MAX = 100      # mini player, search results

# WebP encoding parameters
WEBP_QUALITY = 85     # matches Go's JPEG quality
WEBP_METHOD = 6       # 0=fast 6=slow (slower = better compression)
WEBP_LOSSLESS = False

# Acceptable input MIME types
SUPPORTED_FORMATS = frozenset({"image/jpeg", "image/png", "image/webp", "image/gif"})


# ── Result ─────────────────────────────────────────────────────────────


@dataclass
class CoverOptimizationResult:
    """Three WebP variants of an optimized cover image.

    Each field is ``(image_bytes, width, height)`` or ``None`` if that
    variant was skipped (e.g. original was already ≤ target size).
    """

    original: tuple[bytes, int, int] | None = None  # WebP ≤ 1200px
    med: tuple[bytes, int, int] | None = None        # WebP ≤ 300px
    thumb: tuple[bytes, int, int] | None = None       # WebP ≤ 100px


# ── Errors ─────────────────────────────────────────────────────────────


class CoverOptimizationError(Exception):
    """Raised when image processing fails for any reason."""


# ── Optimizer ──────────────────────────────────────────────────────────


class CoverOptimizer:
    """Optimize cover art images using Pillow.

    Usage::

        optimizer = CoverOptimizer()
        result = optimizer.optimize(image_bytes)
        # result.original[0]  → WebP bytes
        # result.thumb[0]     → 100px WebP bytes
    """

    def optimize(self, data: bytes) -> CoverOptimizationResult:
        """Process raw image bytes into optimized WebP variants.

        Args:
            data: Raw image bytes (JPEG, PNG, WebP, or GIF).

        Returns:
            A ``CoverOptimizationResult`` with up to three WebP variants.

        Raises:
            CoverOptimizationError: If the image cannot be decoded or encoded.
        """
        try:
            img = Image.open(io.BytesIO(data))
        except Exception as exc:
            raise CoverOptimizationError(f"Failed to decode image: {exc}") from exc

        # Convert RGBA/PA/P to RGB for WebP encoding
        if img.mode in ("RGBA", "P", "PA"):
            img = img.convert("RGBA")
        elif img.mode not in ("RGB", "L"):
            img = img.convert("RGB")

        original_w, original_h = img.size

        # Strip EXIF and orientation
        # Pillow's ImageOps.exif_transpose handles orientation tags
        img = ImageOps.exif_transpose(img) or img

        return CoverOptimizationResult(
            original=self._encode_webp(img, ORIGINAL_MAX),
            med=self._encode_webp(img, MED_MAX),
            thumb=self._encode_webp(img, THUMB_MAX),
        )

    # ── Internal ────────────────────────────────────────────────────────

    def _encode_webp(
        self,
        img: Image.Image,
        max_size: int,
    ) -> tuple[bytes, int, int] | None:
        """Resize (preserving aspect ratio) and encode as WebP.

        If the image is already smaller than ``max_size`` on both axes,
        the variant is **not** produced — the caller should fall back to
        the original or next-larger variant.

        Returns:
            ``(webp_bytes, width, height)`` or ``None`` if skipped.
        """
        w, h = img.size

        # Skip if both dimensions are already within limit
        # (except for ORIGINAL_MAX which always produces a result)
        if max_size != ORIGINAL_MAX and w <= max_size and h <= max_size:
            return None

        # Resize preserving aspect ratio
        if w > max_size or h > max_size:
            ratio = min(max_size / w, max_size / h)
            new_w = int(w * ratio)
            new_h = int(h * ratio)
            # Lanczos = high-quality downscale
            resized = img.resize((new_w, new_h), Image.LANCZOS)
        else:
            resized = img

        buf = io.BytesIO()
        try:
            resized.save(
                buf,
                format="WEBP",
                quality=WEBP_QUALITY,
                method=WEBP_METHOD,
                lossless=WEBP_LOSSLESS,
            )
        except Exception as exc:
            raise CoverOptimizationError(
                f"WebP encoding failed at {max_size}px: {exc}"
            ) from exc

        webp_bytes = buf.getvalue()
        logger.debug(
            "Encoded WebP %dx%d → %dx%d (%d bytes, %.1f KB)",
            w, h,
            resized.width, resized.height,
            len(webp_bytes),
            len(webp_bytes) / 1024,
        )
        return webp_bytes, resized.width, resized.height
