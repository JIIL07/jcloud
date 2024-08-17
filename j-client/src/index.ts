    document.addEventListener("DOMContentLoaded", () => {
    const uploadButton = document.getElementById("uploadButton") as HTMLButtonElement;
    const fileInput = document.getElementById("fileInput") as HTMLInputElement;
    const fileListItems = document.getElementById("fileListItems") as HTMLUListElement;

    uploadButton.addEventListener("click", () => {
        const file = fileInput.files?.[0];
        if (file) {
            uploadFile(file);
        }
    });

    const uploadFile = (file: File) => {
        const formData = new FormData();
        formData.append("file", file);

        fetch("/api/v1/upload", {
            method: "POST",
            body: formData,
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    listFiles();
                } else {
                    alert("Upload failed!");
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("Upload failed!");
            });
    };

    const listFiles = () => {
        fetch("/api/v1/files/get")
            .then(response => response.json())
            .then(data => {
                fileListItems.innerHTML = "";
                data.files.forEach((file: { name: string }) => {
                    const li = document.createElement("li");
                    li.textContent = file.name;
                    fileListItems.appendChild(li);
                });
            })
            .catch(error => {
                console.error("Error:", error);
            });
    };

    listFiles();
});
