## PrimeReact Migration Plan

This document breaks down the phases, affected files, and PrimeReact counterparts required to move the Appointment Master web UI (`web/`) onto PrimeReact widgets and icons while retaining the current UX.

---

### Phase 1 – Audit & Mapping

| Current component (`web/src/components`) | PrimeReact target | Notes / required work |
| --- | --- | --- |
| `Button.tsx` | [`Button`](https://primereact.org/button) | Wrap `Button` to preserve `variant` prop; map `variant="primary"` to `severity="primary"`; pass icons via `icon` prop. |
| `Input.tsx` | [`InputText`](https://primereact.org/inputtext) | Maintain full-width styling with existing CSS vars via styled wrapper or `pt` props. |
| `Textarea.tsx` | [`InputTextarea`](https://primereact.org/inputtextarea) | Ensure auto-resize disabled to match current UX unless explicitly turned on. |
| `Select.tsx` | [`Dropdown`](https://primereact.org/dropdown) | Keep controlled + `react-hook-form` bindings; map `options` shape to `optionLabel`/`optionValue`. |
| `Card.tsx` | [`Card`](https://primereact.org/card) | Move current elevations/borders into custom theme to avoid re-styling in every use. |
| `Badge.tsx` | [`Badge`](https://primereact.org/badge) | Keep color tokens using `severity` or custom class. |
| `ListItem.tsx`, `EmptyState.tsx` | [`Panel`, `Message`, `EmptyState` pattern using `Flex` + icons | Compose from `Card`, `IconField`, `Button`. |
| `Field.tsx` | `Fieldset` + PrimeReact input wrappers | Preserve label, description, error message layout; integrate PrimeReact’s `Message` for errors. |
| `Spinner.tsx` | [`ProgressSpinner`](https://primereact.org/progressspinner) | Replace custom SVG spinner. |
| `Skeleton.tsx` | [`Skeleton`](https://primereact.org/skeleton) | Support shapes and widths used today (`rect`, `line`, `circle`). |
| `CopyButton.tsx` | `Button` + [`Tooltip`](https://primereact.org/tooltip) | Use PrimeReact tooltip for copied state if desired. |
| `ChartFrame.tsx`, `LineChart.tsx`, `Sparkline.tsx`, `ChartFrame.tsx` | [`Chart`](https://primereact.org/chart) | PrimeReact wraps Chart.js; re-create time-series (Line) + sparkline using dataset options; maintain legend layout through PrimeReact legend callbacks. |
| `HandUnderline.tsx`, `SuccessBurst.tsx`, `icons` | `primeicons` SVGs or PrimeReact `Icon` | Replace custom lucide icons with closest primeicons; keep bespoke art if no match. |

Feature-level widgets (`web/src/features/**`) using bespoke markup will need equivalent PrimeReact layouts:

- Forms (`auth`, `appointments`, `bookings`, `banlist`): convert inputs/buttons to PrimeReact, use `InputGroup`, `Password`, `Dropdown`, `Calendar`, `MultiSelect`, `Checkbox`, `ToggleButton`.
- Scheduling widgets (`SlotPicker.tsx`, `AvailabilityCalendar.tsx`): leverage `Calendar`, `DatePicker`, and `PickList`/`SelectButton` for slot selection.
- Analytics/Dashboard charts: swap custom SVG charts for PrimeReact `Chart` Line/Bar/Sparkline combos with shared dataset config.
- Tables/List cards (appointments/bookings lists): consider `DataTable`, `ListBox`, or `Timeline` while mirroring present responsive design.

`lucide-react` icon usage (see `rg 'lucide-react' web/src`) will be mapped to `primeicons` identifiers:

| Lucide icon | PrimeReact/PrimeIcons replacement | Files |
| --- | --- | --- |
| `Mail`, `KeySquare`, `CheckCircle2`, `Lock`, `User`, `Eye`, `EyeOff`, `UserPlus`, `LogIn`, `Hash`, `Users`, `Phone`, `FileText`, `Timer`, `Calendar`, `Clock`, `ChevronLeft`, `ChevronRight`, `Shield`, `Info`, `Type` | Use `pi-envelope`, `pi-key`, `pi-check-circle`, `pi-lock`, `pi-user`, `pi-eye`, `pi-eye-slash`, `pi-user-plus`, `pi-sign-in`, `pi-hashtag`, `pi-users`, `pi-phone`, `pi-file`, `pi-stopwatch`, `pi-calendar`, `pi-clock`, `pi-chevron-left`, `pi-chevron-right`, `pi-shield`, `pi-info-circle`, `pi-list` | Update respective feature pages/components to import PrimeReact `IconField` or plain `<i className="pi pi-..." />`. |

---

### Phase 2 – Global Theming & Setup

1. **PrimeReact theme selection**: Adopt a base theme (e.g., `lara-light-blue`) and align it with existing CSS variables in `web/src/theme`. Extend via `:root` variables or PrimeReact `pt` overrides to keep colors (`--primary`, `--bg-elevated`) consistent.
2. **Global config**: In `web/src/main.tsx`, import PrimeReact CSS, `PrimeReactProvider`, enable ripple if desired (`PrimeReact.ripple = true`), and load `primeicons/primeicons.css`.
3. **Styled-components bridge**: Continue exporting themed tokens via `ThemeProvider`; expose them to PrimeReact via CSS variables to avoid dual theming systems.
4. **Accessibility**: Ensure focus outlines mimic current design by customizing `focusRing` tokens in PrimeReact theme.

---

### Phase 3 – Core Primitive Migration

1. **Create wrapper components** in `web/src/components/primereact` (or reuse existing filenames) that internally render PrimeReact components but expose the existing prop API to avoid churn across features.
2. **Inputs & validation**: Update `Field.tsx` plus any `react-hook-form` connections so that refs and `onChange` semantics match PrimeReact controls (`Controller` wrappers where necessary).
3. **Toast & feedback**: Introduce PrimeReact `Toast` provider if required to standardize notifications.
4. **Icons**: Replace `lucide-react` imports with `primeicons` usage via `<i>` tags or `PrimeIcons` helper constants; update styles to ensure proper sizing/alignment.

Deliverable: All shared components now wrap PrimeReact equivalents while tests/build pass.

---

### Phase 4 – Widget & Feature Updates

1. **Charts & analytics**: Replace `LineChart.tsx`/`Sparkline.tsx` and analytics feature charts with `Chart` (Line, Bar) using data from API; ensure responsive legends and tooltips replicate current UX.
2. **Scheduling widgets**: Rebuild `SlotPicker`, `AvailabilityCalendar`, and `BookingForm` inputs using `Calendar`, `SelectButton`, `Chips`, etc., maintaining business logic in hooks.
3. **Auth & forms**: Swap markup in `LoginPage.tsx`, `SignupPage.tsx`, `VerifyPage.tsx` with PrimeReact form components, ensuring accessible labels and helper text via `Message`.
4. **Lists/cards**: Convert dashboard cards and appointment/bookings cards to use `Card`, `DataView`, or `ListBox`. Respect existing responsive CSS via PrimeFlex utilities or custom CSS modules.
5. **Cross-cutting icons**: Ensure every icon reference is backed by `primeicons` or a custom SVG component if no native alternative exists.

---

### Phase 5 – Validation & Documentation

1. **Regression testing**: Run `npm run build`, `npm run dev` smoke, and manual verification across auth, booking, analytics, and notifications flows.
2. **Cross-browser check**: Verify positioning/styling in Chrome + Firefox (and Safari if available) because PrimeReact components may differ from custom CSS.
3. **Bug tracking**: Record discovered gaps via `bugtrack add` with reproduction steps and fixes.
4. **Documentation**: Update `README.md` or `docs/` with new component usage patterns and any custom PrimeReact theme notes so future work stays consistent.

---

Following these phases will transition the existing styled-component primitives, widgets, and icons to their PrimeReact counterparts while minimizing churn in the feature code and ensuring theming consistency.

