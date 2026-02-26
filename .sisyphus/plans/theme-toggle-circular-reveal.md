# 主题切换圆形揭示动画 - 实施计划

## TL;DR

> 为 PicoClaw Web 前端添加圆形揭示动画主题切换效果，使用 Chrome 原生的 View Transitions API 实现从切换按钮位置展开的视觉动画。
> 
> **Deliverables**:
> - 修改 `ThemeService` 添加 View Transition 封装
> - 修改 `LayoutComponent` 集成坐标传递和统一主题管理
> - 添加全局 CSS 动画样式到 `styles.scss`
> - 移除 LayoutComponent 中重复的主题逻辑
> 
> **Estimated Effort**: Short (~30-60分钟)
> **Parallel Execution**: NO - 顺序依赖
> **Critical Path**: ThemeService → LayoutComponent → 样式优化

---

## Context

### Original Request
用户希望为 PicoClaw Web 前端添加主题切换的圆形揭示动画，要求：
- 只考虑最新 Chrome 兼容性
- 从顶部导航栏的切换按钮位置展开
- 流畅的视觉动画效果

### Interview Summary
**Key Discussions**:
- 项目确认: PicoClaw Web 前端 (Angular + Material)
- 触发位置: 顶部导航栏
- 动画原点: 从切换按钮位置展开
- 技术方案: View Transitions API (Chrome 111+)

### Research Findings
**现有代码结构**:
1. **ThemeService** (`web/src/app/services/theme.service.ts`) - 主题管理信号服务
2. **LayoutComponent** (`web/src/app/components/layout/layout.component.ts`) - 包含主题切换按钮，但有重复逻辑
3. **styles.scss** (`web/src/styles.scss`) - 全局样式，Material 主题配置

**问题识别**:
- LayoutComponent 中存在与 ThemeService 重复的主题切换逻辑
- 需要统一主题管理到 ThemeService

**技术选型**:
- 使用 View Transitions API (Chrome 111+)
- 比纯 CSS 动画性能更好，代码更简洁
- 原生 GPU 加速，无需额外依赖

---

## Work Objectives

### Core Objective
实现一个从切换按钮位置展开的圆形揭示动画，在切换亮色/暗色主题时提供流畅的视觉过渡效果。

### Concrete Deliverables
- `ThemeService.toggleThemeWithTransition(event)` 方法
- `LayoutComponent` 调用新方法并移除重复逻辑
- `styles.scss` 添加圆形揭示动画 CSS
- 降级方案（不支持 API 的浏览器直接切换）

### Definition of Done
- [ ] 点击主题按钮时出现圆形揭示动画
- [ ] 动画从按钮位置平滑展开到全屏
- [ ] 主题正确切换（亮色↔暗色）
- [ ] 不支持 API 的浏览器正常切换（无动画）
- [ ] Playwright 截图验证动画效果

### Must Have
- View Transitions API 集成
- 从按钮位置开始的圆形展开
- 降级方案（直接切换）
- 减少动画偏好支持 (`prefers-reduced-motion`)

### Must NOT Have (Guardrails)
- 不支持旧版浏览器的 polyfill（用户只要求 Chrome）
- 过度复杂的自定义动画库
- 修改其他无关组件

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (Vitest)
- **Automated tests**: NO (此功能以视觉验证为主)
- **Agent-Executed QA**: YES - 使用 Playwright MCP 截图验证

### QA Policy
每个任务包含 Agent-Executed QA Scenarios，验证通过 Playwright MCP 执行并截图。

---

## Execution Strategy

### Sequential Execution (顺序执行)

由于任务间有依赖关系，采用顺序执行而非并行：

```
Wave 1: ThemeService 增强
├── Task 1: 添加 View Transition 支持到 ThemeService
└── Task 2: 添加坐标感知主题切换方法

Wave 2: LayoutComponent 集成
├── Task 3: 修改 LayoutComponent 使用 ThemeService
└── Task 4: 移除重复的主题逻辑

Wave 3: 样式与优化
├── Task 5: 添加全局圆形揭示动画 CSS
└── Task 6: 添加可访问性支持 (prefers-reduced-motion)

Wave 4: 验证
├── Task 7: 使用 Playwright MCP 截图验证动画
└── Task 8: 构建并部署测试

Critical Path: Task 1 → Task 2 → Task 3 → Task 5 → Task 7
```

### Agent Dispatch Summary

所有任务均为 `quick` 类别，涉及少量文件修改。

---

## TODOs

- [ ] 1. 添加 View Transition 支持到 ThemeService

  **What to do**:
  - 在 `ThemeService` 中添加 `toggleThemeWithTransition(event: MouseEvent)` 方法
  - 使用 `document.startViewTransition()` API
  - 传递点击坐标给 CSS 变量 `--theme-toggle-x`, `--theme-toggle-y`
  - 提供降级方案（API 不支持时直接调用 `toggleTheme()`）

  **Must NOT do**:
  - 不要修改现有 `toggleTheme()` 方法的行为
  - 不要删除现有方法，保持向后兼容

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 单一服务修改，逻辑清晰
  - **Skills**: `angular-best-practices`
    - 保持 Angular 服务最佳实践
  - **Skills Evaluated but Omitted**:
    - 不需要 playwright - 此任务纯代码实现

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 2, Task 3
  - **Blocked By**: None

  **References**:
  **Pattern References**:
  - `web/src/app/services/theme.service.ts` - 现有 ThemeService 实现
  - `web/src/app/components/layout/layout.component.ts:215-219` - 当前调用方式

  **External References**:
  - [View Transitions API - Chrome Developers](https://developer.chrome.com/docs/web-platform/view-transitions)
  - [MDN: Document.startViewTransition()](https://developer.mozilla.org/en-US/docs/Web/API/Document/startViewTransition)

  **Acceptance Criteria**:
  - [ ] 新方法 `toggleThemeWithTransition` 已添加到 ThemeService
  - [ ] 方法接受 `MouseEvent` 参数获取点击坐标
  - [ ] 使用 `document.startViewTransition()` 包装主题切换
  - [ ] 设置 CSS 变量 `--theme-toggle-x` 和 `--theme-toggle-y`
  - [ ] 降级方案：API 不存在时调用 `toggleTheme()`

  **QA Scenarios**:
  ```
  Scenario: ThemeService 有 View Transition 方法
    Tool: Bash (cat/grep)
    Preconditions: 文件存在
    Steps:
      1. 运行: grep -n "toggleThemeWithTransition" web/src/app/services/theme.service.ts
      2. 验证: 方法存在且包含 startViewTransition
    Expected Result: 找到方法定义
    Evidence: 终端输出截图
  ```

  **Commit**: YES
  - Message: `feat(theme): add View Transition support to ThemeService`
  - Files: `web/src/app/services/theme.service.ts`

---

- [ ] 2. 验证 ThemeService 增强功能

  **What to do**:
  - 运行 TypeScript 类型检查确保无错误
  - 验证新方法可以正确导入和使用

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 1

  **Acceptance Criteria**:
  - [ ] `cd web && npx tsc --noEmit` 无错误
  - [ ] 新服务可以在组件中注入和使用

  **QA Scenarios**:
  ```
  Scenario: TypeScript 类型检查通过
    Tool: Bash
    Steps:
      1. cd /Users/ijkzen/Projects/GO-Project/picoclaw/web
      2. npx tsc --noEmit
    Expected Result: 无错误输出
    Evidence: 终端输出截图
  ```

  **Commit**: NO (与 Task 1 合并提交)

---

- [ ] 3. 修改 LayoutComponent 使用 ThemeService

  **What to do**:
  - 注入 `ThemeService` 到 `LayoutComponent`
  - 修改模板中的 `(click)` 事件传递 `$event`
  - 创建 `onToggleTheme(event: MouseEvent)` 方法
  - 方法中调用 `themeService.toggleThemeWithTransition(event)`

  **Must NOT do**:
  - 不要立即删除旧的 `toggleTheme()` 方法（留到 Task 4）

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `angular-best-practices`

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 1
  - **Blocks**: Task 4

  **References**:
  **Pattern References**:
  - `web/src/app/components/layout/layout.component.ts:89-94` - 当前按钮实现
  - `web/src/app/components/layout/layout.component.ts:215-219` - 当前 toggleTheme 方法

  **Acceptance Criteria**:
  - [ ] `ThemeService` 已注入到 constructor
  - [ ] 模板按钮 `(click)="onToggleTheme($event)"`
  - [ ] `onToggleTheme` 方法调用 `themeService.toggleThemeWithTransition(event)`

  **QA Scenarios**:
  ```
  Scenario: LayoutComponent 正确调用 ThemeService
    Tool: Bash (grep)
    Steps:
      1. 验证 ThemeService 注入: grep -n "ThemeService" web/src/app/components/layout/layout.component.ts
      2. 验证 onToggleTheme: grep -n "onToggleTheme" web/src/app/components/layout/layout.component.ts
      3. 验证模板修改: grep -n "onToggleTheme" web/src/app/components/layout/layout.component.ts
    Expected Result: 所有检查都找到匹配
    Evidence: 终端输出
  ```

  **Commit**: YES
  - Message: `feat(layout): integrate ThemeService for animated toggle`
  - Files: `web/src/app/components/layout/layout.component.ts`

---

- [ ] 4. 移除 LayoutComponent 重复的主题逻辑

  **What to do**:
  - 删除 `LayoutComponent` 中的 `isDarkMode` signal
  - 删除 `toggleTheme()` 方法
  - 删除 `applyTheme()` 方法
  - 删除 `ngOnInit` 中的主题加载逻辑
  - 模板中使用 `themeService.isDarkMode()` 替代本地 signal

  **Must NOT do**:
  - 不要删除 `ThemeService` 的任何代码

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `angular-best-practices`

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 3

  **Acceptance Criteria**:
  - [ ] LayoutComponent 不再管理主题状态
  - [ ] 模板正确绑定 `themeService.isDarkMode()`
  - [ ] 本地 `isDarkMode` signal 已删除
  - [ ] TypeScript 检查通过

  **QA Scenarios**:
  ```
  Scenario: 重复逻辑已移除
    Tool: Bash
    Steps:
      1. grep -n "private applyTheme" web/src/app/components/layout/layout.component.ts
      2. grep -n "isDarkMode = signal" web/src/app/components/layout/layout.component.ts
    Expected Result: 无匹配（方法已删除）
    Evidence: 终端输出
  ```

  **Commit**: YES
  - Message: `refactor(layout): remove duplicate theme logic, use ThemeService`
  - Files: `web/src/app/components/layout/layout.component.ts`

---

- [ ] 5. 添加全局圆形揭示动画 CSS

  **What to do**:
  - 在 `styles.scss` 中添加 View Transition 动画样式
  - 定义 `::view-transition-new(root)` 和 `::view-transition-old(root)` 样式
  - 使用 `clip-path: circle()` 实现圆形揭示
  - 添加平滑的动画时间函数

  **CSS 实现细节**:
  ```scss
  // 圆形揭示动画 - View Transitions API
  :root {
    --theme-toggle-x: 50%;
    --theme-toggle-y: 50%;
  }

  ::view-transition-new(root) {
    clip-path: circle(0% at var(--theme-toggle-x) var(--theme-toggle-y));
    animation: circular-reveal 0.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
  }

  ::view-transition-old(root) {
    animation: none;
  }

  @keyframes circular-reveal {
    to {
      clip-path: circle(150% at var(--theme-toggle-x) var(--theme-toggle-y));
    }
  }
  ```

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 涉及 CSS 动画和视觉效果
  - **Skills**: `tailwind-design-system`
    - 与现有 Tailwind + Material 样式集成

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 1, Task 3

  **References**:
  **Pattern References**:
  - `web/src/styles.scss` - 现有全局样式
  - `web/src/styles.scss:1-50` - Material theming 设置

  **External References**:
  - [View Transitions - Chrome Developers](https://developer.chrome.com/docs/web-platform/view-transitions/same-document)

  **Acceptance Criteria**:
  - [ ] CSS 变量 `--theme-toggle-x` 和 `--theme-toggle-y` 已定义
  - [ ] `::view-transition-new(root)` 有 clip-path 动画
  - [ ] `circular-reveal` 关键帧定义正确
  - [ ] 动画时间函数平滑 (cubic-bezier)

  **QA Scenarios**:
  ```
  Scenario: CSS 样式已添加
    Tool: Bash (grep)
    Steps:
      1. grep -n "view-transition" web/src/styles.scss
      2. grep -n "circular-reveal" web/src/styles.scss
      3. grep -n "theme-toggle-x" web/src/styles.scss
    Expected Result: 所有样式都存在
    Evidence: 终端输出
  ```

  **Commit**: YES
  - Message: `feat(styles): add circular reveal animation for theme toggle`
  - Files: `web/src/styles.scss`

---

- [ ] 6. 添加可访问性支持 (prefers-reduced-motion)

  **What to do**:
  - 添加 `prefers-reduced-motion` 媒体查询
  - 为偏好减少动画的用户禁用或缩短动画

  **CSS 代码**:
  ```scss
  @media (prefers-reduced-motion: reduce) {
    ::view-transition-old(root),
    ::view-transition-new(root) {
      animation-duration: 0.01s !important;
    }
  }
  ```

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 5

  **Acceptance Criteria**:
  - [ ] `prefers-reduced-motion` 媒体查询已添加
  - [ ] 动画时长被覆盖为极短或禁用

  **Commit**: NO (与 Task 5 合并)

---

- [ ] 7. 使用 Playwright MCP 截图验证动画

  **What to do**:
  - 构建前端 `pnpm run build`
  - 复制到后端 `pkg/web/dist`
  - 构建并安装后端 `make install`
  - 启动 gateway `picoclaw gateway start`
  - 使用 Playwright MCP 打开页面并截图
  - 点击主题切换按钮并录制动画

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `playwright`

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocked By**: Task 5

  **QA Scenarios**:
  ```
  Scenario: 主题切换动画正常工作
    Tool: Playwright MCP
    Preconditions: Gateway 已启动
    Steps:
      1. 导航到 http://127.0.0.1:18791/
      2. 设置视口 1280x800
      3. 截图保存为 initial-state.png
      4. 点击主题切换按钮（选择器: button[matTooltip="Toggle theme"]）
      5. 等待 600ms
      6. 截图保存为 after-toggle.png
      7. 比较两张截图：暗色主题应已应用
      8. 再次点击切换按钮
      9. 等待 600ms
      10. 截图保存为 back-to-light.png
    Expected Result: 
      - after-toggle.png 显示暗色主题
      - back-to-light.png 显示亮色主题
    Evidence: .sisyphus/evidence/theme-toggle-*.png
  ```

  **Commit**: NO (验证步骤)

---

- [ ] 8. 构建并部署测试

  **What to do**:
  - 运行完整构建流程
  - 验证无构建错误
  - 在浏览器中手动测试动画

  **Build Commands**:
  ```bash
  cd /Users/ijkzen/Projects/GO-Project/picoclaw/web && pnpm run build
  rm -rf /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/*
  cp -r /Users/ijkzen/Projects/GO-Project/picoclaw/web/dist/web/* /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/
  cd /Users/ijkzen/Projects/GO-Project/picoclaw && make install
  picoclaw gateway stop && sleep 1 && picoclaw gateway start
  ```

  **Acceptance Criteria**:
  - [ ] 前端构建成功
  - [ ] 后端构建成功
  - [ ] Gateway 正常启动
  - [ ] Web UI 可访问

  **QA Scenarios**:
  ```
  Scenario: 完整构建流程
    Tool: Bash
    Steps:
      1. 执行上述构建命令
      2. 检查每个命令的退出码
      3. 验证 gateway 状态: picoclaw gateway status
    Expected Result: 所有命令成功，gateway 运行中
    Evidence: 终端输出截图
  ```

  **Commit**: NO (构建产物不提交)

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  读取计划检查所有任务已完成：ThemeService 新方法、LayoutComponent 调用、styles.scss 动画、Playwright 截图证据存在。
  Output: `Must Have [4/4] | Evidence [N files] | VERDICT: APPROVE/REJECT`

- [ ] F2. **代码质量检查** — `unspecified-high`
  运行 `cd web && npx tsc --noEmit` 确保无类型错误。检查是否有 `any` 类型滥用或逻辑问题。
  Output: `TypeCheck [PASS/FAIL] | Issues [N] | VERDICT`

- [ ] F3. **视觉验证** — `playwright`
  使用 Playwright 录制主题切换动画视频，验证圆形揭示效果从按钮位置展开。
  Output: `Animation [PASS/FAIL] | Evidence: .sisyphus/evidence/theme-toggle-video.webm`

- [ ] F4. **可访问性检查** — `quick`
  验证 `prefers-reduced-motion` 媒体查询存在且有效。
  Output: `A11y [PASS/FAIL] | VERDICT`

---

## Commit Strategy

| Task | Commit Message | Files |
|------|---------------|-------|
| 1-2 | `feat(theme): add View Transition support to ThemeService` | `web/src/app/services/theme.service.ts` |
| 3 | `feat(layout): integrate ThemeService for animated toggle` | `web/src/app/components/layout/layout.component.ts` |
| 4 | `refactor(layout): remove duplicate theme logic` | `web/src/app/components/layout/layout.component.ts` |
| 5-6 | `feat(styles): add circular reveal animation for theme toggle` | `web/src/styles.scss` |

---

## Success Criteria

### Verification Commands
```bash
# TypeScript 类型检查
cd web && npx tsc --noEmit

# 前端构建
cd web && pnpm run build

# 后端构建
cd .. && make install

# Gateway 状态
picoclaw gateway status

# 截图验证（Playwright MCP 执行）
# - 导航到 http://127.0.0.1:18791/
# - 点击主题切换按钮
# - 验证暗色主题应用
```

### Final Checklist
- [ ] ThemeService 有 `toggleThemeWithTransition` 方法
- [ ] LayoutComponent 使用 ThemeService 并移除重复逻辑
- [ ] styles.scss 有圆形揭示动画 CSS
- [ ] 降级方案工作正常（直接切换）
- [ ] prefers-reduced-motion 支持已添加
- [ ] Playwright 截图验证动画效果
- [ ] TypeScript 类型检查通过
- [ ] 构建成功，Gateway 运行正常

---

## 实现要点速查

### ThemeService 新增方法
```typescript
async toggleThemeWithTransition(event?: MouseEvent): Promise<void> {
  if (!('startViewTransition' in document)) {
    this.toggleTheme();
    return;
  }

  const x = event?.clientX ?? window.innerWidth / 2;
  const y = event?.clientY ?? window.innerHeight / 2;

  document.documentElement.style.setProperty('--theme-toggle-x', `${x}px`);
  document.documentElement.style.setProperty('--theme-toggle-y', `${y}px`);

  const transition = document.startViewTransition(() => {
    this.toggleTheme();
  });

  await transition.ready;
}
```

### styles.scss 新增样式
```scss
:root {
  --theme-toggle-x: 50%;
  --theme-toggle-y: 50%;
}

::view-transition-new(root) {
  clip-path: circle(0% at var(--theme-toggle-x) var(--theme-toggle-y));
  animation: circular-reveal 0.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

::view-transition-old(root) {
  animation: none;
}

@keyframes circular-reveal {
  to {
    clip-path: circle(150% at var(--theme-toggle-x) var(--theme-toggle-y));
  }
}

@media (prefers-reduced-motion: reduce) {
  ::view-transition-old(root),
  ::view-transition-new(root) {
    animation-duration: 0.01s !important;
  }
}
```

### LayoutComponent 修改
```typescript
// constructor 注入
constructor(
  private breakpointObserver: BreakpointObserver,
  private themeService: ThemeService
) {}

// 新方法
onToggleTheme(event: MouseEvent): void {
  this.themeService.toggleThemeWithTransition(event);
}
```

```html
<!-- 模板修改 -->
<button
  mat-icon-button
  (click)="onToggleTheme($event)"
  matTooltip="Toggle theme">
  <mat-icon>{{ themeService.isDarkMode() ? 'light_mode' : 'dark_mode' }}</mat-icon>
</button>
```

---

**Plan saved to**: `.sisyphus/plans/theme-toggle-circular-reveal.md`
