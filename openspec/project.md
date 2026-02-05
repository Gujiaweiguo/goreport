# Project Context

## Purpose
JimuReport is a data visualization reporting tool with an online designer for
reports, printing, dashboards, and large-screen displays. This repository
includes a Spring Boot example app under jimureport-example for integration and
demo use.

## Tech Stack
- Java 17
- Spring Boot 3.5.x
- Maven
- JimuReport/JimuBI Spring Boot starters
- MySQL 5.7+ (required for the example app)
- Redis (optional)
- Sa-Token (auth)
- Freemarker (templates)
- Lombok
- Druid (DB pool)
- JSQLParser
- Minio (optional)

## Project Conventions

### Code Style
- 4-space indentation; opening brace on the same line
- No wildcard imports; keep imports explicit
- SLF4J logging (often via Lombok @Slf4j) with parameterized messages
- Mixed Chinese/English comments and messages; preserve local style
- Avoid adding new dependencies without confirmation

### Architecture Patterns
- Base package: com.jeecg (modules under com.jeecg.modules.jmreport)
- Spring MVC controllers in controller, config in config, extensions in extend
- Sa-Token utilities under satoken
- Field injection via @Autowired is common
- Controllers often return view names or redirect: strings

### Testing Strategy
- Tests are minimal; use Maven Surefire when needed
- Common commands: mvn test, mvn -Dtest=TestClassName test
- No explicit lint/format tooling detected; follow existing style

### Git Workflow
- Not documented; keep changes small and focused

## Domain Context
- JimuReport focuses on complex reports and printing; JimuBI focuses on
  dashboards and large-screen visualizations
- Example app endpoints: /jmreport/list and /drag/list
- Default example credentials: admin/123456
- Supports many SQL/NoSQL data sources (MySQL, Oracle, SQLServer, Postgres, etc.)

## Important Constraints
- Requires JDK 17+ for the Spring Boot 3 example app
- Requires MySQL 5.7+; initialize with db/jimureport.mysql5.7.create.sql
- Redis integration is optional (dependencies commented out)
- License: LGPL; commercial license required to remove branding per README

## External Dependencies
- MySQL database (required), Redis (optional)
- Optional Minio integration
- Maven repositories: Aliyun and Jeecg
- Official docs and quick start at https://help.jimureport.com
