# TODO

## High Priority
1. **App Menus Refactor:** Replace verbose `MenuItem::Custom` blocks with `non_updateable_custom_item()` helpers in `app/src/app_menus.rs` to clean up codebase.
2. **Launch Config IDs:** Add proper Launch Config IDs to data sources.
3. **Agent Mode cleanup:** Remove deprecated AI assistant codebase elements that conflict with Agent Mode.

## Medium Priority
4. **Search Bar Placeholders:** Implement real support for placeholder text in the search bar.
5. **UI Rendering:** Address layout issues in `ActionButton` when children exceed maximum widths.
6. **Viewporting:** Fix viewported elements inside the welcome and command palettes.

## Low Priority
7. **Native App Menus:** Implement native Windows OS app menus.
8. **Font Caching:** Implement retry mechanisms for partially loaded font families.
