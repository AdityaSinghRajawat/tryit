// Phase 1: only Swagger UI is wired up here. Phase 2 plugs in Redoc, Mintlify,
// and a generic heuristic detector.

import * as swagger from "./frameworks/swagger";

export type Framework = "swagger" | "none";

export function detect(): Framework {
  if (swagger.isSwaggerUI()) return "swagger";
  return "none";
}

export const frameworks = { swagger };
