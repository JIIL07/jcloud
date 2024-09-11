use eframe::egui;
use rfd::FileDialog;
use serde::{Deserialize, Serialize};
use std::fs::{self, File};
use std::io::Write;
use std::path::{Path, PathBuf};

#[derive(Serialize, Deserialize, Debug, Clone)]
struct FileItem {
    name: String,
    path: String,
    selected: bool,
    is_dir: bool,
    depth: usize,
}

impl FileItem {
    fn new(path: PathBuf, is_dir: bool, depth: usize) -> Self {
        let name = path
            .file_name()
            .unwrap_or_else(|| path.as_os_str())
            .to_string_lossy()
            .to_string();

        FileItem {
            name,
            path: path.to_string_lossy().to_string(),
            selected: false,
            is_dir,
            depth,
        }
    }
}

struct MyApp {
    root_dir: Option<PathBuf>,
    file_tree: Vec<FileItem>,
    show_save_button: bool,
}

impl Default for MyApp {
    fn default() -> Self {
        Self {
            root_dir: None,
            file_tree: vec![],
            show_save_button: false,
        }
    }
}

impl MyApp {
    fn build_file_tree(&mut self, dir: &Path) {
        self.file_tree.clear();
        self.visit_dirs(dir, 0);
    }

    fn visit_dirs(&mut self, dir: &Path, depth: usize) {
        if dir.is_dir() {
            self.file_tree.push(FileItem::new(dir.to_path_buf(), true, depth));

            let mut dirs = vec![];
            let mut files = vec![];

            for entry in fs::read_dir(dir).unwrap() {
                let entry = entry.unwrap();
                let path = entry.path();
                if path.is_dir() {
                    dirs.push(path);
                } else {
                    files.push(path);
                }
            }

            for dir in dirs {
                self.visit_dirs(&dir, depth + 1);
            }

            for file in files {
                self.file_tree.push(FileItem::new(file, false, depth + 1));
            }
        }
    }

    fn save_selection(&self) {
        let selected_items: Vec<FileItem> = self
            .file_tree
            .iter()
            .filter(|item| item.selected)
            .cloned()
            .collect();

        let json = serde_json::to_string_pretty(&selected_items).unwrap();
        let mut file = File::create("selected.json").unwrap();
        file.write_all(json.as_bytes()).unwrap();
    }

    fn update_selection(&mut self, idx: usize, selected: bool) {
        self.file_tree[idx].selected = selected;
        if self.file_tree[idx].is_dir {
            let depth = self.file_tree[idx].depth;
            for i in idx + 1..self.file_tree.len() {
                if self.file_tree[i].depth > depth {
                    self.file_tree[i].selected = selected;
                } else {
                    break;
                }
            }
        }
    }
}

impl eframe::App for MyApp {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Jcloud interactive");

            if ui.button("Pick directory").clicked() {
                if let Some(path) = FileDialog::new().pick_folder() {
                    self.root_dir = Some(path.clone());
                    self.build_file_tree(&path);
                    self.show_save_button = true;
                }
            }

            if let Some(root) = &self.root_dir {
                ui.label(format!("Selected directory: {}", root.display()));

                let mut changes = vec![];

                ui.separator();

                if self.show_save_button {
                    if ui.button("Save selection").clicked() {
                        self.save_selection();
                        self.show_save_button = false;
                    }
                }

                egui::ScrollArea::vertical()
                    .auto_shrink([false; 2])
                    .show(ui, |ui| {
                        ui.set_min_height(ui.available_height());
                        ui.set_min_width(ui.available_width());

                        for (idx, item) in self.file_tree.iter_mut().enumerate() {
                            ui.horizontal(|ui| {
                                ui.add_space(item.depth as f32 * 20.0);
                                let label = if item.is_dir {
                                    format!("📁 {}", item.name)
                                } else {
                                    format!("📄 {}", item.name)
                                };
                                let selected = ui.checkbox(&mut item.selected, label.as_str()).changed();
                                if selected {
                                    changes.push((idx, item.selected));
                                }
                            });
                        }
                    });

                for (idx, selected) in changes {
                    self.update_selection(idx, selected);
                }
            }
        });
    }
}

fn main() {
    let options = eframe::NativeOptions {
        initial_window_size: Some(egui::vec2(400.0, 500.0)),
        resizable: true,
        ..Default::default()
    };
    eframe::run_native("jcloud", options, Box::new(|_cc| Box::new(MyApp::default())))
}
