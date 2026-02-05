<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# AGENTS.md
# Guidance for agentic coding in this repository.

## Scope and layout
- Primary runnable project lives in `jimureport-example/`.
- Main Spring Boot entrypoint in code: `com.jeecg.JimuReportApplication`.
- Docs reference running from `jimureport-example` and Docker usage.

## Build, lint, test, run

### Maven (jimureport-example)
- Build package (documented):
  - `mvn clean package`
  - Source: `jimureport-example/README.md`
- Run app (documented as IDE run):
  - Run main class from IDE.
  - README mentions `org.jeecg.modules.JimuReportApplication`, but actual class is `com.jeecg.JimuReportApplication`.
  - Sources: `jimureport-example/README.md`, `jimureport-example/src/main/java/com/jeecg/JimuReportApplication.java`.
- Test (standard Maven, not documented in repo):
  - `mvn test`
  - `mvn -Dtest=TestClassName test`
  - `mvn -Dtest=TestClassName#testMethod test`

### Docker (jimureport-example)
- Build image stack (documented):
  - `docker-compose up -d`
  - Source: `jimureport-example/README.md`

### Database bootstrap
- Initialize DB schema (documented):
  - Execute `db/jimureport.mysql5.7.create.sql`.
  - Source: `jimureport-example/README.md`

### Lint/format
- No explicit lint/format tooling found (no Checkstyle/PMD/SpotBugs/EditorConfig).
- Do not assume auto-formatters exist; follow existing style by observation.

## Code style guidelines (observed)

### Language and frameworks
- Java 17 (see `jimureport-example/pom.xml`).
- Spring Boot 3.5.x with standard Spring annotations.
- Lombok used for logging (`@Slf4j`).

### Project structure
- Base package: `com.jeecg` / `com.jeecg.modules.jmreport`.
- Common package areas:
  - `controller` for MVC controllers.
  - `config` for configuration classes.
  - `extend` for JimuReport extension implementations.
  - `satoken` for auth/session utilities.

### Imports
- No wildcard imports observed.
- Imports grouped by third-party first, then JDK.
- Keep imports minimal and explicit.

### Formatting and layout
- 4-space indentation.
- Opening brace on the same line.
- Blank line between logical sections (imports, class fields, methods).
- Inline `if` spacing with braces:
  - `if (condition) {` and `} else {`.

### Naming conventions
- Classes: `PascalCase` (e.g., `LoginController`).
- Methods/fields: `camelCase` (e.g., `getToken`).
- Constants: `UPPER_SNAKE_CASE` (e.g., `LOGIN_PAGE`).
- Packages: lowercase dotted names.
- YAML keys: lowercase with hyphens (e.g., `token-name`).

### Dependency injection
- Field injection via `@Autowired` is used (no constructor injection observed).
- Fields are often package-private (no access modifier).

### Logging
- Uses SLF4J (`LoggerFactory`) or Lombok `@Slf4j`.
- Log levels:
  - `info` for normal flow.
  - `warn` for recoverable issues.
  - `error` for failures.
- Prefer parameterized logging (`{}`) over string concatenation.

### Error handling
- Prefer explicit exception handling with logging.
- Custom runtime exception used: `JimuReportException`.
- Exception handler classes exist but are commented out.
- Avoid empty catch blocks; if swallowing, add a comment or log.

### HTTP and web conventions
- Controllers use `@Controller`, `@GetMapping`, `@RequestMapping`.
- Redirects constructed via `redirect:` prefix.
- Response bodies are mostly string views, not JSON in controllers.

### Auth/session conventions (Sa-Token)
- Sa-Token used for auth (`cn.dev33.satoken.*`).
- Token is read from Sa-Token context or from request parameter `token`.
- Token name configured in `application.yml` (`X-Access-Token`).

### Comments and docs
- JavaDoc used for public methods in some classes.
- Inline comments used for non-obvious behavior (e.g., redirects, token handling).

## Testing guidance
- Test sources are minimal; do not assume extensive coverage.
- Use Maven Surefire single-test runs when possible.
- If adding tests, align package structure with `src/main/java`.
- There is a sample data class under `src/main/java` (`testdb`), not in `src/test`.
- Do not move or rename it unless asked; treat it as runtime sample data.

## Environment assumptions
- JDK 17+ required.
- MySQL 5.7+ required for example app.
- Redis optional and controlled via optional dependencies.

## Docker notes
- Docker setup documented under `jimureport-example/README.md`.
- Mac M-series guidance requires ARM base image tweaks (see README).

## Localization
- Comments and log messages include both Chinese and English.
- Preserve existing language style in nearby code; do not rewrite for consistency.

## Repo-specific notes
- `.gitignore` ignores `target/`, `logs/`, and IDE files.
- No CI configuration found in `.github/workflows` or other CI files.

## Cursor/Copilot rules
- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` found.

## 规则
- 默认中文回复
- 有歧义先提问；先给计划/方案，得到确认再执行
- 最小改动，不做无关重构
- 不新增依赖，需先说明并确认
- 未明确要求不提交；不使用 --amend

## OpenSpec（如使用）
- 新功能/改接口/改表结构：先 Proposal，经批准后 Apply
- 上线后 Archive，并执行 `openspec validate --strict --no-interactive`

## Sources consulted
- `jimureport-example/README.md`
- `jimureport-example/README.en-US.md`
- `jimureport-example/pom.xml`
- `jimureport-example/src/main/java/com/jeecg/JimuReportApplication.java`
- `jimureport-example/src/main/java/com/jeecg/modules/jmreport/controller/LoginController.java`
- `jimureport-example/src/main/java/com/jeecg/modules/jmreport/extend/JimuReportTokenServiceImpl.java`
- `jimureport-example/src/main/java/com/jeecg/modules/jmreport/satoken/util/AjaxRequestUtils.java`
- `jimureport-example/src/main/java/com/jeecg/modules/jmreport/config/CustomCorsConfiguration.java`
- `jimureport-example/src/main/resources/application.yml`
- `.gitignore`
