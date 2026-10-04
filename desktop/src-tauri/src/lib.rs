#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    // The webview has no shell or filesystem permissions. Cluster credentials
    // stay in the Go API process and are never handed to this shell.
    tauri::Builder::default()
        .run(tauri::generate_context!())
        .expect("error while running KubeMv");
}
