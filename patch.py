with open("./crates/computer_use/src/lib.rs", "r") as f:
    content = f.read()

new_content = content.replace("""    pub fn validate(&self) -> Result<(), String> {
        if self.top_left.x() < 0 || self.top_left.y() < 0 {
            return Err(format!(
                "Screenshot region top_left must be non-negative, got ({}, {})",
                self.top_left.x(),
                self.top_left.y()
            ));
        }
        if self.bottom_right.x() <= self.top_left.x() {""", """    pub fn validate(&self) -> Result<(), String> {
        if self.bottom_right.x() <= self.top_left.x() {""")

with open("./crates/computer_use/src/lib.rs", "w") as f:
    f.write(new_content)
