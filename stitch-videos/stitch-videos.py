#!/usr/bin/env python3
"""
Video Stitching Tool
====================

PURPOSE:
    Extracts time-based snippets from multiple video sources (local files or YouTube URLs)
    and stitches them together into a single output video. Supports audio normalization
    and fade effects.

USAGE:
    ./stitch-videos.py <config.yaml>

    OR

    python3 stitch-videos.py <config.yaml>

REQUIREMENTS:
    - Python 3.7+
    - FFmpeg installed and available in PATH
    - Python libraries: yt-dlp, moviepy (2.x), PyYAML

    Install dependencies:
        pip install yt-dlp moviepy PyYAML

YAML CONFIGURATION FILE:
    The script requires a YAML configuration file with the following structure:

    output_file: (string, optional)
        Name of the final output video file.
        Default: "final_video.mp4"

    normalize_audio: (boolean, optional)
        Whether to normalize audio levels across all clips.
        Default: true

    output_fps: (number, optional)
        Frame rate (frames per second) for the output video.
        Helps ensure smooth playback when combining videos with different frame rates.
        Common values: 24, 30, 60
        Default: 30

    output_height: (number, optional)
        Output video height in pixels. Width is automatically scaled to maintain aspect ratio.
        Lower resolution = MUCH faster processing (especially for 4K sources).
        Common values: 480 (SD), 720 (HD), 1080 (Full HD), 2160 (4K)
        Default: None (keeps source resolution)

    encoding_preset: (string, optional)
        FFmpeg encoding preset that controls speed vs quality trade-off.
        Options: ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow
        Faster presets = quicker encoding but larger files
        Slower presets = better compression but much slower
        Default: medium

    video_bitrate: (string, optional)
        Video bitrate for output quality (e.g., "2000k", "5000k", "10000k").
        Higher bitrate = better quality but larger file size
        Common values: 2000k (low), 5000k (good), 8000k (high), 15000k (very high)
        Default: 5000k

    thread_count: (number, optional)
        Number of CPU threads to use for encoding.
        More threads = faster encoding (up to your CPU core count)
        Set to 0 for auto-detection (uses all available cores)
        Default: 4

    video_tasks: (list, required)
        List of video sources and their snippets to extract.

        Each task contains:
            location: (string, required)
                Path to local video file OR YouTube URL (any yt-dlp supported URL)

            snippets: (list, required)
                List of time-based segments to extract from this video.

                Each snippet contains:
                    start: (string or number, required)
                        Start time in seconds (90) or timestamp format (1:30 or 1:30:45)
                        Supports fractional seconds: 10.5, 1:30.5, 1:30:45.6

                    end: (string or number, required)
                        End time in seconds or timestamp format
                        Supports fractional seconds: 10.5, 1:30.5, 1:30:45.6

                    fade_in: (number, optional)
                        Fade-in duration in seconds
                        Default: 0 (no fade)

                    fade_out: (number, optional)
                        Fade-out duration in seconds
                        Default: 0 (no fade)

EXAMPLE YAML:
    ---
    output_file: "my_compilation.mp4"
    normalize_audio: true
    output_fps: 30
    output_height: 1080
    encoding_preset: medium
    video_bitrate: 5000k
    thread_count: 4

    video_tasks:
    - location: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
      snippets:
      - start: "0:30"
        end: "1:45"
        fade_in: 1.0
        fade_out: 1.0

      - start: "2:15.5"
        end: "3:00.25"
        fade_out: 0.5

    - location: "/path/to/local/video.mp4"
      snippets:
      - start: 10
        end: 25
        fade_in: 0.5

      - start: 45
        end: "1:30"

    - location: "https://youtube.com/shorts/abc123xyz"
      snippets:
      - start: 0
        end: 15

PERFORMANCE TUNING:
    Video processing is CPU/GPU intensive. These settings control speed vs quality:

    ENCODING PRESET (encoding_preset):
        Controls the speed/quality/filesize trade-off:

        ultrafast   - 5-10x faster than medium, largest files, slightly lower quality
        superfast   - 4-7x faster than medium
        veryfast    - 3-5x faster than medium (RECOMMENDED for speed)
        faster      - 2-3x faster than medium
        fast        - 1.5-2x faster than medium
        medium      - Balanced (DEFAULT) - good for most use cases
        slow        - Better compression, 2x slower
        slower      - Even better compression, 3-4x slower
        veryslow    - Best compression, 5-10x slower

        💡 TIP: Use 'veryfast' for quick previews, 'medium' for final output

    VIDEO BITRATE (video_bitrate):
        Higher bitrate = better quality + larger files + slightly slower encoding

        "2000k"     - Low quality (web previews, low-res videos)
        "5000k"     - Good quality (DEFAULT) - suitable for 1080p
        "8000k"     - High quality (1080p with lots of motion)
        "15000k"    - Very high quality (4K or professional work)

        💡 TIP: Lower bitrate significantly speeds up encoding

    FRAME RATE (output_fps):
        Lower FPS = fewer frames to process = faster

        24          - Cinematic look, 20% faster than 30 FPS
        30          - Standard web video (DEFAULT)
        60          - Smooth motion for gaming/sports, 2x slower than 30 FPS

        💡 TIP: If source videos are 24 FPS, output at 24 FPS for best speed

    OUTPUT RESOLUTION (output_height):
        ⚡ BIGGEST PERFORMANCE IMPACT - Resize video to lower resolution

        Processing 4K video (2160p) is EXTREMELY slow. Scaling down dramatically speeds up:

        None        - Keep source resolution (DEFAULT) - slowest for 4K sources
        1080        - Full HD - 4x faster than 4K
        720         - HD - 9x faster than 4K
        480         - SD - 20x faster than 4K

        💡 TIP: If your source is 4K, use output_height: 1080 for 4x speedup
        💡 TIP: For quick previews, use 720 or even 480
        💡 TIP: Width is automatically calculated to preserve aspect ratio

    THREAD COUNT (thread_count):
        More threads = faster encoding (up to your CPU core count)

        1-8         - Specific thread count
        0           - Auto-detect (uses all available CPU cores)
        Default: 4

        💡 TIP: Run 'nproc' (Linux) or 'sysctl -n hw.ncpu' (Mac) to see your cores
        💡 TIP: Set to 0 or your core count for maximum speed

    AUDIO NORMALIZATION (normalize_audio):
        Audio normalization adds processing time

        false       - Skip normalization, faster processing
        true        - Normalize audio (DEFAULT), ensures consistent volume

        💡 TIP: Disable if all source videos have similar volume levels

    EXAMPLE PERFORMANCE CONFIGURATIONS:

        Fast Preview (10-20x faster, especially for 4K sources):
            output_height: 720
            encoding_preset: veryfast
            video_bitrate: 2000k
            output_fps: 24
            thread_count: 0
            normalize_audio: false

        Balanced (DEFAULT):
            encoding_preset: medium
            video_bitrate: 5000k
            output_fps: 30
            thread_count: 4
            normalize_audio: true

        High Quality 1080p from 4K source (4x faster than keeping 4K):
            output_height: 1080
            encoding_preset: medium
            video_bitrate: 8000k
            output_fps: 30
            thread_count: 0
            normalize_audio: true

        High Quality keeping original resolution (slower):
            encoding_preset: slow
            video_bitrate: 8000k
            output_fps: 60
            thread_count: 0
            normalize_audio: true

NOTES:
- Remote videos are downloaded temporarily and deleted after processing
- All clips are concatenated in the order specified in the YAML file
- The script validates all sources before starting to avoid wasting time
- Large videos or long processing times may require patience
"""

import os
import sys
import shutil
import subprocess

# --- PRE-FLIGHT CHECKS ---
missing_libs = []
try:
    import yt_dlp
except ImportError:
    missing_libs.append("yt-dlp")
try:
    from moviepy import VideoFileClip, concatenate_videoclips
    from moviepy import afx, vfx
except ImportError:
    missing_libs.append("moviepy")
try:
    import yaml
except ImportError:
    missing_libs.append("PyYAML")

if missing_libs:
    print(f"\n[!] MISSING LIBRARIES: pip install {' '.join(missing_libs)}")
    sys.exit(1)

if shutil.which("ffmpeg") is None:
    print("\n[!] FFmpeg NOT FOUND. Please install it to proceed.")
    sys.exit(1)

# --- UTILITIES ---

def timestamp_to_seconds(ts):
    if isinstance(ts, (int, float)):
        return ts
    parts = ts.split(':')
    if len(parts) == 1:
        return float(parts[0])
    elif len(parts) == 2:
        return int(parts[0]) * 60 + float(parts[1])
    elif len(parts) == 3:
        return int(parts[0]) * 3600 + int(parts[1]) * 60 + float(parts[2])
    return ts

def is_remote_valid(url):
    """Checks if a YouTube URL is valid and accessible without downloading."""
    ydl_opts = {'quiet': True, 'no_warnings': True, 'simulate': True}
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        try:
            ydl.extract_info(url, download=False)
            return True
        except Exception:
            return False

def download_video(url, output_path):
    ydl_opts = {
        'format': 'bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best',
        'outtmpl': output_path,
        'quiet': True,
        'no_warnings': True,
    }
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])

# --- MAIN PROCESS ---

def main():
    if len(sys.argv) < 2:
        print("\nUsage: python script.py <your_config.yaml>")
        sys.exit(1)

    with open(sys.argv[1], 'r') as f:
        config = yaml.safe_load(f)

    tasks = config.get('video_tasks', [])
    output_name = config.get('output_file', 'final_video.mp4')
    normalize = config.get('normalize_audio', True)
    output_fps = config.get('output_fps', 30)
    output_height = config.get('output_height', None)
    encoding_preset = config.get('encoding_preset', 'medium')
    video_bitrate = config.get('video_bitrate', '5000k')
    thread_count = config.get('thread_count', 4)
    temp_dir = "temp_processing"

    # --- CONFIGURATION SUMMARY ---
    print("--- Configuration ---")
    print(f"Output File: {output_name}")
    print(f"Normalize Audio: {normalize}")
    print(f"Output FPS: {output_fps}")
    print(f"Output Height: {output_height if output_height else 'Original (no scaling)'}")
    print(f"Encoding Preset: {encoding_preset}")
    print(f"Video Bitrate: {video_bitrate}")
    print(f"Thread Count: {thread_count}")
    print()

    # --- FAIL-FAST VALIDATION PHASE ---
    print("--- Validating All Input Sources ---")
    invalid_sources = []

    for task in tasks:
        loc = task.get('location', '')
        if loc.lower().startswith("http"):
            print(f" Checking Remote: {loc}...", end="\r")
            if not is_remote_valid(loc):
                invalid_sources.append(f"Remote (Invalid/Private): {loc}")
        else:
            print(f" Checking Local: {loc}...", end="\r")
            if not os.path.exists(loc):
                invalid_sources.append(f"Local (File Not Found): {loc}")

    if invalid_sources:
        print("\n\n[!] CRITICAL ERROR: Validation failed. Process aborted.")
        for err in invalid_sources:
            print(f"    - {err}")
        sys.exit(1)

    print("\nAll sources verified. Starting execution...\n")

    # --- EXECUTION PHASE ---
    if not os.path.exists(temp_dir):
        os.makedirs(temp_dir)

    all_snippets = []
    video_objects = []

    try:
        for i, task in enumerate(tasks):
            location = task['location']
            snippets = task.get('snippets', [])
            is_url = location.lower().startswith("http")

            if is_url:
                source_path = os.path.join(temp_dir, f"remote_source_{i}.mp4")
                print(f"Processing Scene {i+1} (Remote Download): {location}")
                download_video(location, source_path)
            else:
                source_path = location
                print(f"Processing Scene {i+1} (Local File): {source_path}")

            video = VideoFileClip(source_path)
            video_objects.append(video)

            for j, s in enumerate(snippets):
                start = timestamp_to_seconds(s['start'])
                end = timestamp_to_seconds(s['end'])

                print(f"  > Snippet {j+1}: {s['start']} to {s['end']}")
                clip = video.subclipped(start, end)

                if output_height:
                    clip = clip.with_effects([vfx.Resize(height=output_height)])

                effects = []
                effects.append(vfx.EvenSize())
                if normalize and clip.audio is not None:
                    effects.append(afx.AudioNormalize())
                if s.get('fade_in', 0) > 0:
                    effects.append(vfx.FadeIn(s['fade_in']))
                if s.get('fade_out', 0) > 0:
                    effects.append(vfx.FadeOut(s['fade_out']))

                clip = clip.with_effects(effects)

                all_snippets.append(clip)

        if all_snippets:
            print("\n--- Final Rendering ---")
            final = concatenate_videoclips(all_snippets, method="compose")
            final.write_videofile(
                output_name,
                codec="libx264",
                audio_codec="aac",
                fps=output_fps,
                preset=encoding_preset,
                bitrate=video_bitrate,
                threads=thread_count,
                ffmpeg_params=["-pix_fmt", "yuv420p", "-profile:v", "high"]
            )
            print(f"\nSUCCESS: Created {output_name}")
        else:
            print("\n[!] No clips were generated.")

    finally:
        print("\nCleaning up temporary files...")
        for v in video_objects:
            v.close()
        if os.path.exists(temp_dir):
            shutil.rmtree(temp_dir)

if __name__ == "__main__":
    main()
