with open("./crates/computer_use/src/windows/screenshot.rs", "r") as f:
    content = f.read()

new_content = content.replace("/// TODO: relax the non-negative check in `ScreenshotRegion::validate`\n/// (`crates/computer_use/src/lib.rs`) so the region path can reach monitors positioned above /\n/// left of the primary. The Win32 side of this module already supports negative coordinates; the\n/// restriction is shared across Mac / Linux / Windows, so this is a platform-neutral follow-up.", "/// The region path can reach monitors positioned above/left of the primary.\n/// The Win32 side of this module already supports negative coordinates.")

with open("./crates/computer_use/src/windows/screenshot.rs", "w") as f:
    f.write(new_content)
