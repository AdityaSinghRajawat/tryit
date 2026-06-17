// Framework detection. Order of precedence: swagger (deterministic DOM,
// produces a client-side hint) → redoc (deterministic DOM, server cascade)
// → generic (heading+codeblock heuristic, server cascade).

import * as swagger from "./frameworks/swagger";
import * as redoc from "./frameworks/redoc";
import * as generic from "./frameworks/generic";

export type Framework = "swagger" | "redoc" | "generic" | "none";

export function detect(): Framework {
  if (swagger.isSwaggerUI()) return "swagger";
  if (redoc.isRedoc()) return "redoc";
  if (generic.findEndpointBlocks().length > 0) return "generic";
  return "none";
}

export const frameworks = { swagger, redoc, generic };
