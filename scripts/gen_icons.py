#!/usr/bin/env python3
"""Generate the KubeMv desktop icons."""

import struct
import zlib
from pathlib import Path

OUT = Path(__file__).resolve().parents[1] / "desktop" / "src-tauri" / "icons"


def chunk(tag: bytes, data: bytes) -> bytes:
    return struct.pack(">I", len(data)) + tag + data + struct.pack(">I", zlib.crc32(tag + data) & 0xFFFFFFFF)


def write_png(path: Path, size: int) -> bytes:
    raw = bytearray()
    for y in range(size):
        raw.append(0)
        for x in range(size):
            raw.extend(pixel(x, y, size))
    ihdr = struct.pack(">IIBBBBB", size, size, 8, 6, 0, 0, 0)
    png = b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", ihdr) + chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + chunk(b"IEND", b"")
    path.write_bytes(png)
    return png


def pixel(x: int, y: int, size: int) -> tuple[int, int, int, int]:
    nx = (x + 0.5) / size
    ny = (y + 0.5) / size
    if not inside_round_rect(nx, ny, 0.16):
        return (0, 0, 0, 0)
    background = (16, 20, 26, 255)
    brass = (215, 161, 95, 255)
    if 0.22 <= nx <= 0.36 and 0.18 <= ny <= 0.82:
        return brass
    if near_segment(nx, ny, 0.34, 0.50, 0.78, 0.20, 0.055):
        return brass
    if near_segment(nx, ny, 0.34, 0.50, 0.78, 0.80, 0.055):
        return brass
    return background


def inside_round_rect(nx: float, ny: float, radius: float) -> bool:
    inset = 0.06
    left, top, right, bottom = inset, inset, 1 - inset, 1 - inset
    cx = min(max(nx, left + radius), right - radius)
    cy = min(max(ny, top + radius), bottom - radius)
    return (nx - cx) ** 2 + (ny - cy) ** 2 <= radius ** 2


def near_segment(px: float, py: float, ax: float, ay: float, bx: float, by: float, thickness: float) -> bool:
    abx, aby = bx - ax, by - ay
    length_sq = abx * abx + aby * aby
    if length_sq == 0:
        return False
    t = max(0.0, min(1.0, ((px - ax) * abx + (py - ay) * aby) / length_sq))
    dx, dy = px - (ax + t * abx), py - (ay + t * aby)
    return dx * dx + dy * dy <= thickness * thickness


def write_ico(path: Path, png: bytes, size: int) -> None:
    width = 0 if size >= 256 else size
    header = struct.pack("<HHH", 0, 1, 1)
    entry = struct.pack("<BBBBHHII", width, width, 0, 0, 1, 32, len(png), 22)
    path.write_bytes(header + entry + png)


def write_icns(path: Path, png: bytes) -> None:
    inner = b"ic07" + struct.pack(">I", len(png) + 8) + png
    path.write_bytes(b"icns" + struct.pack(">I", len(inner) + 8) + inner)


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    write_png(OUT / "32x32.png", 32)
    png128 = write_png(OUT / "128x128.png", 128)
    write_png(OUT / "128x128@2x.png", 256)
    write_ico(OUT / "icon.ico", png128, 128)
    write_icns(OUT / "icon.icns", png128)


if __name__ == "__main__":
    main()
