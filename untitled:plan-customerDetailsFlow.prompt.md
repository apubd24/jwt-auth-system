Build a complete customer management feature for an existing Go + React project with JWT auth and admin/read-only roles.

### Backend requirements

- Use Go with Gin, GORM, and PostgreSQL.
- Add a `Customer` model with:
  - `company_name`
  - `office_address`
  - `website`
  - `support_email`
  - `noc_email`
  - `logo_url`
  - `status` (`active` / `inactive`)
  - `created_at`, `updated_at`

- Add related models:
  - `CommonPhone` for multiple company-wide phone/WhatsApp numbers
  - `ContactPerson` with:
    - `name`
    - `designation`
    - `contact_type` (`technical`, `sales`, `level-1`, etc.)
    - multiple contact items: email, mobile, whatsapp
  - `ContactMethod` or nested arrays to support multiple emails/mobiles/WhatsApp per contact person

- Migrate models with GORM and `OnDelete:CASCADE` relationships.

- Add APIs:
  - `GET /api/customers` — view all customers
  - `GET /api/customers/:id` — view a single customer
  - `POST /api/customers` — create customer
  - `PUT /api/customers/:id` — update customer
  - `DELETE /api/customers/:id` — delete customer
  - `PUT /api/customers/:id/status` — activate/deactivate customer

- Enforce roles:
  - admin only: create / update / delete / status change
  - readonly: view only GET endpoints
- Use existing middleware style with `AuthMiddleware()` and `RequireRole("admin")`
- Return JSON responses with proper status codes and validation errors

### Frontend requirements

- Build a React page using Tailwind CSS
- Use a token from `localStorage.token` and `localStorage.role`
- Create:
  - customer create/edit form
  - dynamic “Add Contact Person” button
  - dynamic contact-person fields for name, designation, contact type, email, mobile, WhatsApp
  - dynamic company phone fields
  - list of customers with edit/delete actions
- admin can create/edit/delete customers
- readonly can only view list and details
- keep simple Tailwind form style consistent with an existing Tailwind project
- use Axios for API calls

### Deliverables

- Go model files for customer data
- Go handler file(s) for customer CRUD and status toggling
- Router updates for customer endpoints
- React component(s) for customer management page
- API helper file for customer requests
- Tailwind styling for dynamic form fields and table/list

### Context

- Current backend package path is `backend/`
- Use `jwt-auth-backend` module structure
- Use existing JWT auth middleware and role middleware setup
- Frontend is a React app using Tailwind CSS

---

Use that prompt to generate the full backend and frontend implementation.