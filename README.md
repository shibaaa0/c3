# 📂 c3 - Terminal File Manager

**c3** là một trình quản lý file nhẹ chạy trong terminal, được viết bằng **Go** với thư viện [tview](https://github.com/rivo/tview).
Ứng dụng cung cấp một **giao diện TUI đơn giản, đầy màu sắc** để duyệt và quản lý file/thư mục một cách hiệu quả.

![Main Interface](main_interface.png)

---

## ✨ Tính năng

- 🔽 **Điều hướng**: Duyệt thư mục bằng phím mũi tên, ấn `S` để hiện/ẩn file ẩn.  
- 🔍 **Tìm kiếm**: Lọc file/thư mục theo tiền tố bằng phím `F`.  
- 🖼 **Preview Pane**: Xem nội dung thư mục hoặc metadata file *(size, thời gian sửa đổi, quyền)*.  
- 🎨 **UI nhiều màu**:
  - 📘 Thư mục: xanh dương  
  - 📄 File: trắng  
  - ✅ Thành công: xanh lá  
  - ❌ Lỗi: đỏ  

![Feature Preview](feature_preview.png)

---


## ⚡ Cài đặt

### 🔑 Yêu cầu
- [Go](https://go.dev/dl/): Phiên bản **1.16+**  
  _(kiểm tra bằng `go version`)_

### 🚀 Các bước
```bash
# Clone repository
git clone https://github.com/shibaaa0/c3.git
cd c3/app/

# Cài dependencies
go mod tidy

# Chạy trực tiếp
go run main.go

# Hoặc build & cài đặt
go build -o c3
sudo mv c3 /usr/local/bin/
c3
```

---

## 🎮 Sử dụng

Chạy lệnh:

```bash
c3
```

### ⌨️ Key Bindings

| Phím | Hành động |
|------|-----------|
| ←    | Lùi về thư mục cha |
| s    | Hiện/ẩn file ẩn |
| f    | Tìm kiếm file/thư mục |
| Esc  | Thoát ứng dụng |

![Key Bindings](key_bindings.png)

---

![Demo Workflow](demo_workflow.gif)
