use druid::{AppLauncher, Color, Data, Lens, Widget, WidgetExt, WindowDesc};
use druid::widget::{Button, Flex, TextBox};
use std::fs::File;
use std::io::prelude::*;

#[derive(Clone, Data, Lens)]
struct AppState {
    file_content: String,
}

fn main() {
    let main_window = WindowDesc::new(build_ui())
        .title("Druid File Upload Example")
        .window_size((600.0, 400.0));

    let initial_state = AppState { file_content: String::new() };

    AppLauncher::with_window(main_window)
        .log_to_console()
        .launch(initial_state)
        .expect("Failed to launch application");
}

fn build_ui() -> impl Widget<AppState> {
    let text_box = TextBox::new()
        .with_placeholder("Enter file content")
        .padding(10.0)
        .lens(AppState::file_content);

    let button = Button::new("Save File")
        .on_click(|_ctx, data: &mut AppState, _env| {
            let mut file = File::create("uploaded_file.txt").expect("Unable to create file");
            file.write_all(data.file_content.as_bytes()).expect("Unable to write data");
        })
        .padding(10.0);

    Flex::column()
        .with_child(text_box)
        .with_child(button)
        .background(Color::from_hex_str("#E1E9F0").unwrap())
        .center()
}
