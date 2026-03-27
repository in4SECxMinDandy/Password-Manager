# Frontend Architecture

## 1. Overview

React + TypeScript frontend with Vite build system, Tailwind CSS styling, and Shadcn/ui components.

## 2. Technology Stack

| Category | Technology |
|----------|------------|
| Framework | React 18 |
| Language | TypeScript |
| Build | Vite |
| Styling | Tailwind CSS |
| Components | Shadcn/ui (Radix primitives) |
| Routing | React Router 6 |
| State | React Context + TanStack Query |
| Forms | React Hook Form + Zod |
| Notifications | Sonner (toast) |
| Icons | Lucide React |

## 3. Directory Structure

```
frontend/
├── src/
│   ├── app/                    # Route pages
│   │   ├── (auth)/             # Auth routes
│   │   │   ├── login/
│   │   │   └── register/
│   │   └── (app)/              # Protected routes
│   │       ├── vaults/
│   │       └── entries/
│   │
│   ├── components/
│   │   ├── ui/                 # Shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   └── dropdown-menu.tsx
│   │   │
│   │   └── layout/             # Layout components
│   │       ├── AppLayout.tsx
│   │       └── Sidebar.tsx
│   │
│   ├── features/
│   │   ├── auth/
│   │   │   ├── components/      # Login, Register pages
│   │   │   ├── context/         # AuthContext
│   │   │   └── hooks/          # useAuth hook
│   │   │
│   │   └── vault/
│   │       ├── components/     # Vault, Entry pages
│   │       └── hooks/          # useVault hook
│   │
│   └── lib/
│       ├── api-client.ts       # API client
│       ├── api-types.ts        # TypeScript types
│       └── utils.ts            # Utility functions
```

## 4. Design System

### Color Palette

Uses CSS variables for theming:

```css
:root {
  --background: hsl(0 0% 100%);
  --foreground: hsl(222.2 84% 4.9%);
  --primary: hsl(222.2 47.4% 11.2%);
  --primary-foreground: hsl(210 40% 98%);
  /* ... */
}

.dark {
  --background: hsl(222.2 84% 4.9%);
  --foreground: hsl(210 40% 98%);
  /* ... */
}
```

### Typography

System font stack with Tailwind:
- Font: System UI, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto

### Spacing

Tailwind default spacing scale (0.25rem increments)

## 5. Key Components

### Auth Components

| Component | Description |
|-----------|-------------|
| `LoginPage` | Email/password form with show/hide password |
| `RegisterPage` | Registration with password strength indicator |
| `AuthContext` | Authentication state management |

### Vault Components

| Component | Description |
|-----------|-------------|
| `VaultsPage` | Grid of vaults with create/delete |
| `EntriesPage` | List of entries with search/filter |
| `Sidebar` | Navigation with user menu, theme toggle |

## 6. API Client

The `api-client.ts` provides:

- Automatic token refresh on 401
- Type-safe API methods
- Error handling with typed errors

```typescript
// Usage
const { data } = useQuery({
  queryKey: ['vaults'],
  queryFn: () => apiClient.getVaults(),
})
```

## 7. Routing

```
/login              - Login page (public)
/register          - Register page (public)
/vaults             - Vault list (protected)
/vaults/:id/entries - Entry list (protected)
```

Protected routes wrap with `AuthProvider` and check authentication.

## 8. Roadmap

- [ ] Password generator component
- [ ] MFA setup UI
- [ ] Entry detail modal
- [ ] Keyboard shortcuts
- [ ] Mobile responsive refinements
