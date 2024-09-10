use eframe::egui;
use rfd::FileDialog;
use serde::{Deserialize, Serialize};
use std::fs::{self, File};
use std::io::Write;
use std::path::{Path, PathBuf};

// Структура для хранения информации о выбранных файлах
#[derive(Serialize, Deserialize, Debug, Clone)]
struct FileItem {
    path: String,
    selected: bool,
}

impl FileItem {
    fn new(path: String) -> Self {
        FileItem {
            path,
            selected: false,
        }
    }
}

// Основная структура приложения
struct MyApp {
    root_dir: Option<PathBuf>,
    file_tree: Vec<FileItem>,
}

impl Default for MyApp {
    fn default() -> Self {
        Self {
            root_dir: None,
            file_tree: vec![],
        }
    }
}

impl MyApp {
    // Функция для построения дерева файлов
    fn build_file_tree(&mut self, dir: &Path) {
        self.file_tree.clear();
        self.visit_dirs(dir);
    }

    // Рекурсивная функция обхода директорий
    fn visit_dirs(&mut self, dir: &Path) {
        if dir.is_dir() {
            for entry in fs::read_dir(dir).unwrap() {
                let entry = entry.unwrap();
                let path = entry.path();
                if path.is_dir() {
                    self.visit_dirs(&path);
                }
                self.file_tree.push(FileItem::new(path.to_string_lossy().to_string()));
            }
        }
    }

    // Сохранение выбранных файлов в JSON
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
}

// Реализация пользовательского интерфейса
impl eframe::App for MyApp {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Файловый менеджер");

            // Кнопка для выбора директории
            if ui.button("Выбрать директорию").clicked() {
                if let Some(path) = FileDialog::new().pick_folder() {
                    self.root_dir = Some(path.clone());
                    self.build_file_tree(&path);
                }
            }

            // Если директория выбрана, строим дерево файлов
            if let Some(root) = &self.root_dir {
                ui.label(format!("Выбранная директория: {}", root.display()));

                egui::ScrollArea::vertical().show(ui, |ui| {
                    for item in &mut self.file_tree {
                        ui.horizontal(|ui| {
                            ui.checkbox(&mut item.selected, item.path.as_str());
                        });
                    }
                });
            }

            // Кнопка для сохранения выбора
            if ui.button("Save Selection").clicked() {
                self.save_selection();
            }
        });
    }
}

fn main() {
    let options = eframe::NativeOptions {
        initial_window_size: Some(egui::vec2(600.0, 400.0)),
        ..Default::default()
    };
    eframe::run_native("Файловый менеджер", options, Box::new(|_cc| Box::new(MyApp::default())))
}
