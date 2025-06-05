# lkvm

**lkvm** is a lightweight, open-source USB-based KVM (Keyboard-Video-Mouse) over IP solution, enabling remote control of machines via USB interfaces.

---

## Features

* **USB-Based KVM**: Utilizes USB interfaces to provide KVM functionality over IP.

* **Cross-Platform Compatibility**: Supports multiple operating systems for both host and client.

* **Minimal Dependencies**: Designed with a focus on simplicity and minimal external dependencies.

* **Open-Source**: Released under the MIT license, encouraging community contributions and transparency.

---

## Installation

To install **lkvm**, clone the repository and build the project using Go:

```bash
git clone https://github.com/tshelter/lkvm.git
cd lkvm
go build ./cmd/lkvm
```


Ensure that Go is installed and properly configured on your system.

---

## Usage

After building, run the **lkvm** binary to start the KVM service:

```bash
./lkvm
```

---

# 🚀 Open-Source Projects

## 🐍 Python Library for CH9329

A library for working with the CH9329 chip in Python. Originally developed by others, but extensively debugged and improved.

[🔗 GitHub Repository](https://github.com/tshelter/py-ch9329/)

## ⚡ Minimal PoC Prototype

A minimal proof-of-concept implementation.

[🔗 GitHub Repository](https://github.com/tshelter/pylkvm)

## 🏗️ Go Library for CH9329

A Go library for working with the CH9329 chip, built from scratch with a complete redesign of the approach.

[🔗 GitHub Repository](https://github.com/tshelter/ch9329/)

## 🔌 Fully Functional USB-KVM Solution

A complete working solution for portable USB-KVM emulation.

[🔗 GitHub Repository](https://github.com/tshelter/lkvm)
