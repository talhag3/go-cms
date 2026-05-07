┌─────────────────────────────────────────────────────────────────────────┐
│                            HTTP REQUEST                                 │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         MIDDLEWARE LAYER                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │ Recover  │→│ Logger   │→│RequestID │→│ CORS     │               │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘               │
│                                                      │                  │
│                                              ┌───────▼───────┐         │
│                                              │  Auth (if    │         │
│                                              │  protected)  │         │
│                                              └───────┬───────┘         │
└──────────────────────────────────────────────────────┼──────────────────┘
                                                       │
                                                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         HANDLER LAYER (Controllers)                     │
│  ┌─────────────────────┐  ┌─────────────────────┐                      │
│  │   PostHandler       │  │   UserHandler       │                      │
│  │  - GetAll()         │  │  - GetAll()         │                      │
│  │  - GetByID()        │  │  - GetByID()        │                      │
│  │  - Create()         │  │  - Create()         │                      │
│  │  - Update()         │  │  - Login()          │                      │
│  │  - Delete()         │  │  - Delete()         │                      │
│  └──────────┬──────────┘  └──────────┬──────────┘                      │
└─────────────┼────────────────────────┼──────────────────────────────────┘
              │                        │
              ▼                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          SERVICE LAYER (Business Logic)                 │
│  ┌─────────────────────┐  ┌─────────────────────┐                      │
│  │   PostService       │  │   UserService       │                      │
│  │  - Validation       │  │  - Authentication   │                      │
│  │  - Business Rules   │  │  - User Logic       │                      │
│  │  - Enrichment       │  │                     │                      │
│  └──────────┬──────────┘  └──────────┬──────────┘                      │
└─────────────┼────────────────────────┼──────────────────────────────────┘
              │                        │
              ▼                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        REPOSITORY LAYER (Data Access)                   │
│  ┌─────────────────────┐  ┌─────────────────────┐                      │
│  │  PostRepository     │  │  UserRepository     │                      │
│  │  (Interface)        │  │  (Interface)        │                      │
│  └──────────┬──────────┘  └──────────┬──────────┘                      │
│             │                        │                                  │
│             ▼                        ▼                                  │
│  ┌─────────────────────┐  ┌─────────────────────┐                      │
│  │MockPostRepository   │  │MockUserRepository   │                      │
│  │  (In-Memory)        │  │  (In-Memory)        │                      │
│  └─────────────────────┘  └─────────────────────┘                      │
│                                                                      │
│  ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─  │
│  PRODUCTION: Replace mock with:                                        │
│  ┌─────────────────────┐  ┌─────────────────────┐                      │
│  │PostgreSQLRepository │  │PostgreSQLRepository │                      │
│  │  (using pgx)        │  │  (using pgx)        │                      │
│  └─────────────────────┘  └─────────────────────┘                      │
└─────────────────────────────────────────────────────────────────────────┘


---

## Best Practices Summary

┌────────────────────────────────────────────────────────────────────────┐
│                     GO BEST PRACTICES FOR PHP DEVS                     │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  1. FOLDER = NAMESPACE                                                 │
│     PHP: namespace App\Services;                                       │
│     Go:  package services (folder name = package name)                 │
│                                                                        │
│  2. CAPITALIZATION = VISIBILITY                                        │
│     PHP: public/private/protected                                      │
│     Go:  Capital = exported, lowercase = unexported                    │
│                                                                        │
│  3. ERROR HANDLING                                                      │
│     PHP: try/catch with exceptions                                     │
│     Go:  if err != nil { return err } (explicit, no hidden control)    │
│                                                                        │
│  4. DEPENDENCY INJECTION                                                │
│     PHP: Constructor injection via container                           │
│     Go:  Pass interfaces via constructor functions                     │
│                                                                        │
│  5. INTERFACES FOR ABSTRACTION                                          │
│     PHP: Explicit "implements"                                         │
│     Go:  Implicit (if it has the methods, it implements)               │
│     WHY: Makes swapping implementations easy (mock → real DB)          │
│                                                                        │
│  6. KEEP MAIN() THIN                                                    │
│     Only: config loading, DI wiring, server start                      │
│     All logic in packages                                              │
│                                                                        │
│  7. INTERNAL/ FOR PRIVATE CODE                                          │
│     Code in internal/ can't be imported by other projects              │
│     Like "private" packages                                            │
│                                                                        │
│  8. PKG/ FOR REUSABLE CODE                                             │
│     Code in pkg/ is designed to be reused                              │
│     Like "public" packages                                             │
│                                                                        │
│  9. USE POINTERS WISELY                                                 │
│     Use *T when:                                                        │
│       - Need nil (null)                                                │
│       - Need to modify original                                        │
│       - Avoid copying large structs                                    │
│     Don't use when: simple values, primitives                          │
│                                                                        │
│  10. DEFER FOR CLEANUP                                                  │
│      PHP: try/finally or __destruct                                    │
│      Go:  defer func() { ... } (executes when function returns)        │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘