 Run a quick tool diagnostic: list project files, search for "func", read main.go, explore
  the codebase structure with the Explore agent, create a test todo, search the web for
  "golang testing", and ask me one question. Keep each step minimal.

  This will trigger:
  - Glob - file listing
  - Grep - code search
  - Read - file reading
  - Task (Explore) - codebase exploration agent
  - TodoWrite - task management
  - WebSearch - web search capability
  - AskUserQuestion - user interaction

  For a more comprehensive test including other agents, you could add:

  Run a tool diagnostic: list files, grep for "main", read main.go, use Explore agent for
  codebase structure, use codebase-locator to find test files, create a test todo list, web
  search "golang best practices", and ask me a test question. Be brief.

  This adds Task (codebase-locator) subagent type.
