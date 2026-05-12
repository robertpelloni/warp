# Project TODOs

```
./crates/websocket/src/proxy.rs:6://! TODO: Switch to tungstenite's native proxy support once it is available and remove this
./crates/computer_use/src/lib.rs:217:    // TODO(AGENT-2283): consider making this a type that is cheap to clone
./crates/computer_use/src/linux/x11/screenshot.rs:28:    // TODO: Consider compositing the cursor into the screenshot in the future.
./crates/computer_use/src/mac/keycode_cache.rs:66:/// TODO(QUALITY-271): Store the modifier keys as well.
./crates/syntax_tree/src/queries/indent_query.rs:122:            // TODO(INT-614): Remove this special case.
./crates/syntax_tree/src/queries/highlight_query.rs:97:// TODO(kevin): Once we migrate buffer to store ArrayStrings. We should implement the chunks API on buffer directly to avoid collecting
./crates/http_server/src/lib.rs:17:    /// TODO(vorporeal): Remove this when we have a shared tokio runtime.
./crates/settings/src/macros.rs:657:                // TODO(advait): deprecate this in favour of struct settings in a follow-up PR.
./crates/warpui_core/src/debug/view_tree_debug_view.rs:54:                // TODO(vorporeal): Don't arbitrarily pick font family 0.
./crates/warpui_core/src/windowing/mod.rs:101:// TODO(CORE-2691): implement native Windows OS app menus
./crates/warpui_core/src/integration/step.rs:644:        // TODO would be cool to move it under IntegrationTestEvent
./crates/warpui_core/src/core/app.rs:901:    /// TODO(CORE-2323): Implement native Windows OS modal
./crates/warpui_core/src/core/app.rs:2638:            // TODO(vorporeal): what's the right value here?
./crates/warpui_core/src/core/app.rs:2642:            // TODO(alokedesai): Determine if, and how, we want to pass the on_gpu_driver_reported
./crates/warpui_core/src/core/app.rs:3385:                    // TODO(PLAT-746): Update API for loading a font family to
./crates/warpui_core/src/core/app.rs:3403:        // TODO(PLAT-760): Implement a retry mechanism if we load some but not
./crates/warpui_core/src/core/app.rs:3426:            // TODO(PLAT-747): Ideally the the font cache itself is a model and can
./crates/warpui_core/src/core/app.rs:3715:        // TODO: Apply the same deferred unsubscribe pattern used in `emit_event` to support
./crates/warpui_core/src/core/app.rs:4367:    // TODO(vorporeal): why is AppContext.window_bounds holding an option?
./crates/warpui_core/src/core/entity.rs:37:/// TODO(vorporeal): This can probably be eliminated entirely, with View and
./crates/warpui_core/src/core/view/context.rs:503:    /// TODO(vorporeal): Determine how best to eliminate this function and move
./crates/warpui_core/src/ui_components/radio_buttons.rs:100:// TODO(roland): Remembering the selected option can be unintuitive if the number of options
./crates/warpui_core/src/ui_components/link.rs:12:    text: String, // TODO figure out how it can be ui element (or icon?)
./crates/warpui_core/src/ui_components/text.rs:171:    // TODO(alokedesai): Make it clear throughout the text rendering code that highlights are
./crates/warpui_core/src/ui_components/components.rs:58:    pub width: Option<f32>, // TODO should be possible to spec units/equations (eg. 100% - 5px)
./crates/warpui_core/src/platform/app.rs:106:    // TODO(CORE-2322): implement desktop notifications on Windows
./crates/warpui_core/src/platform/app.rs:160:    // TODO(CORE-2683): implement events for internet reachability changes
./crates/warpui_core/src/platform/app.rs:295:// TODO(CORE-2691): implement native Windows OS app menus
./crates/warpui_core/src/platform/app.rs:317:// TODO(CORE-2323): implement native Windows OS modal
./crates/warpui_core/src/platform/file_picker.rs:105:    // TODO(CORE-2324): open file picker on Windows
./crates/warpui_core/src/platform/mod.rs:81:// TODO(advait): revisit this to check if there's a better approach.
./crates/warpui_core/src/event.rs:90:/// TODO: for the events that have modifiers (e.g. cmd, shift), we should
./crates/warpui_core/src/image_cache.rs:370:            ImageType::Svg { .. } => 0, // TODO: How do we calculate svg size in bytes?
./crates/warpui_core/src/image_cache.rs:464:    // TODO: other types
./crates/warpui_core/src/image_cache.rs:830:    /// TODO(APP-3877): remove `#[allow(dead_code)]` once the debounce eviction pass wires this up.
./crates/warpui_core/src/scene.rs:442:        // TODO: Investigate how / when we would pass a z-index that isn't in the scene
./crates/warpui_core/src/scene.rs:657:        // TODO: Support hit testing on glyphs?
./crates/warpui_core/src/telemetry/event_store.rs:20:    // TODO: write to disk periodically
./crates/warpui_core/src/text_layout.rs:49:// TODO: Ideally, we would use DEFAULT_MONOSPACE_FONT_SIZE here, however that
./crates/warpui_core/src/elements/text.rs:1208:                // TODO (roland): if the mouse is over a character, position_for_point could put us to the left or
./crates/warpui_core/src/elements/resizable.rs:28:/// TODO:
./crates/warpui_core/src/elements/new_scrollable/mod.rs:36:// TODO: we might want this to be configurable.
./crates/warpui_core/src/elements/new_scrollable/mod.rs:56:/// TODO: currently, this constant reflects the value that makes sense for MacOS (cocoa) scroll events.
./crates/warpui_core/src/elements/hoverable.rs:507:        // TODO: we should re-consider this approach. It can lead to missed `on_hover` dispatches.
./crates/warpui_core/src/elements/clipped_scrollable.rs:183:/// TODO: there is currently a bug with constraint-passing when nesting
./crates/warpui_core/src/elements/formatted_text_element.rs:83:// TODO: We should think about whether line height applies to notebooks as well.
./crates/warpui_core/src/elements/formatted_text_element.rs:288:    /// TODO (roland): There are cases where we need to set the line height to the UI default
./crates/warpui_core/src/elements/formatted_text_element.rs:355:                        // Default hover handler does nothing. TODO: highlighting
./crates/warpui_core/src/elements/formatted_text_element.rs:441:                        // Default hover handler does nothing. TODO: highlighting
./crates/warpui_core/src/elements/formatted_text_element.rs:1633:                // TODO: Update when we support task lists.
./crates/warpui_core/src/elements/formatted_text_element.rs:1659:                    // TODO: In the future, can use the first parameter (the lang field) to modify the actual style of the text fragment.
./crates/warpui_core/src/elements/formatted_text_element.rs:1839:                    // TODO: ensure the constructed text is the same as `line.raw_text()` (test?)
./crates/warpui_core/src/elements/scrollable.rs:41:/// TODO: currently, this constant reflects the value that makes sense for MacOS (cocoa) scroll events.
./crates/warpui_core/src/elements/scrollable.rs:311:        // TODO(kevin): Do we need the scroll_start <= 0 check?
./crates/warpui_core/src/fonts.rs:273:    // TODO(alokedesai): Better consolidate the caching logic between the FontCache and the
./crates/markdown_parser/src/html_parser.rs:277:                            // TODO: Support Github's code block representation.
./crates/markdown_parser/src/html_parser.rs:448:                    // TODO: We need to add more phrasing styling we support (e.g. links) here.
./crates/markdown_parser/src/html_parser_test.rs:191:// TODO: remove/update this test when we eventually support these HTML element types!
./crates/markdown_parser/src/markdown_parser.rs:92:        // TODO: most markdown parsers stop treating this as a list item after some number of spaces > 4
./crates/markdown_parser/src/markdown_parser.rs:1856:    // TODO: Look into other autolink rules here: https://github.github.com/gfm/#autolinks-extension-
./crates/warp_util/src/path.rs:82:    // TODO(peter): we ought to tolerate non-Unicode paths here.
./crates/warp_util/src/path.rs:99:                        // TODO While checking `cfg!(windows)` is usually correct for determining
./crates/warp_util/src/path.rs:308:    // TODO(CORE-2805): Make sure this works in distribution.
./crates/warpui/src/windowing/winit/text_layout_tests.rs:351:// TODO(PLAT-779): check all line bounds once bidirectional wrapping is fixed in cosmic-text.
./crates/warpui/src/windowing/winit/text_layout_tests.rs:430:// TODO(PLAT-779): check all line bounds once bidirectional wrapping is fixed in cosmic-text.
./crates/warpui/src/windowing/winit/text_layout_tests.rs:493:// TODO(PLAT-779): check all line bounds once bidirectional wrapping is fixed in cosmic-text.
./crates/warpui/src/windowing/winit/app.rs:79:    /// TODO(CORE-2274): theming on Windows
./crates/warpui/src/windowing/winit/event_loop/key_events.rs:168:    // TODO(wasm): Extend this to support all of the function/shift/arrow keys.
./crates/warpui/src/windowing/winit/event_loop/key_events_tests.rs:7:    // TODO: it would be nice to test the following:
./crates/warpui/src/windowing/winit/event_loop/mod.rs:277:        // TODO(advait): Implement the function key for winit.
./crates/warpui/src/windowing/winit/event_loop/mod.rs:865:                // TODO(vorporeal): Should we be calling approve_termination() here?
./crates/warpui/src/windowing/winit/event_loop/mod.rs:1162:                        // TODO: when we need key codes for voice input on Linux/Windows, we'll need to populate this!
./crates/warpui/src/windowing/winit/event_loop/mod.rs:1791:            // TODO(abhishek): We make sure that the position is different than last time to prevent winit from
./crates/warpui/src/windowing/winit/delegate.rs:323:        // TODO(wasm): Investigate implementing this by creating a <input> element
./crates/warpui/src/windowing/winit/delegate.rs:497:        // TODO(wasm): implement this.
./crates/warpui/src/windowing/winit/delegate.rs:501:        // TODO(wasm): implement this.
./crates/warpui/src/windowing/winit/delegate.rs:506:        // TODO(wasm): Implement this.
./crates/warpui/src/windowing/winit/delegate.rs:510:        // TODO(wasm): Implement this.
./crates/warpui/src/windowing/winit/delegate.rs:531:        // TODO(wasm): Implement this.
./crates/warpui/src/windowing/winit/delegate.rs:557:        // TODO
./crates/warpui/src/windowing/winit/window.rs:73:        /// TODO(CORE-1891) Instead of being hard-coded, this should be configurable by the user via
./crates/warpui/src/windowing/winit/window.rs:1416:                    // TODO: use location to actually draw buttons
./crates/warpui/src/windowing/winit/linux/window_manager.rs:15:    // TODO(CORE-3034): Re-enable this codepath once we've understood and
./crates/warpui/src/windowing/winit/linux/window_manager.rs:61:/// TODO(CORE-3034): Re-enable this codepath and remove the allow(dead_code)
./crates/warpui/src/windowing/winit/fonts/swash_rasterizer.rs:107:        // TODO(alokedesai): Ensure our font rasterization code is robust to returned formats that
./crates/warpui/src/windowing/winit/fonts.rs:272:                        // TODO(alokedesai): Make refactors to fontdb and/or font-kit to make
./crates/warpui/src/windowing/winit/fonts.rs:485:        // TODO(alokedesai): Consider using FontDB's `make_shared_font_data` here. FontDB creates a temporary memory
./crates/warpui/src/windowing/winit/fonts.rs:601:            // TODO(alokedesai): Properly clip multi-line text using the same strategy we use on mac.
./crates/warpui/src/windowing/winit/fonts.rs:649:            // TODO(daprahamian): when we have time, we should investigate pulling
./crates/warpui/src/windowing/winit/fonts.rs:694:            // TODO(vorporeal): See if we need to compute this (and if so, how to).
./crates/warpui/src/platform/mac/rendering/metal/renderer.rs:996:        // TODO(alokedesai): Backport the optimization to only set the size of surface when a
./crates/warpui/src/platform/mac/event.rs:184:        // TODO: This option is deprecated by Apple in favour of NSEventTypeOtherMouseDown
./crates/warpui/src/platform/mac/window.rs:517:                // TODO: device appears to be leaked here.
./crates/warpui/src/platform/mac/text_layout.rs:836:        // TODO(CORE-2004): If we want to support external font fallback on
```
