# Contributing to SingUI

Thank you for your interest in contributing to **SingUI**! We welcome contributions from developers of all skill levels.

---

## 🛠️ Development Setup

### 1. Prerequisites
- **Go**: 1.22+
- **Node.js**: 20+ & npm / pnpm
- **Sing-box**: 1.9+ (Optional for local proxy testing)

### 2. Fork & Clone
```bash
git clone https://github.com/<your-username>/SingUI.git
cd SingUI
```

### 3. Running Locally

#### Frontend Development Server:
```bash
cd frontend
npm install
npm run dev
# Running on http://localhost:3000 (proxies API calls to port 2096)
```

#### Backend Development Server:
```bash
cd backend
go run ./cmd/server/main.go -p 2096
# Running on http://localhost:2096
```

---

## 📋 Contribution Workflow

1. **Create an Issue**: Before submitting major features, please open an Issue to discuss the architecture and design.
2. **Create a Branch**:
   ```bash
   git checkout -b feat/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```
3. **Write Clean Code & Test**:
   - Follow standard Go idioms (`gofmt`, `golint`).
   - Follow Vue 3 Composition API & TypeScript best practices.
   - Verify frontend build with `npm run build`.
4. **Commit Guidelines**:
   We recommend [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat: add new protocol support`
   - `fix: resolve subscription format parsing error`
   - `docs: update documentation`
   - `refactor: optimize core process supervisor`
5. **Submit a Pull Request**:
   - Open a PR against the `main` branch.
   - Describe the changes and reference any related issues.

---

## 📜 Code of Conduct

Please be respectful, constructive, and collaborative in all communications and discussions within the community.
