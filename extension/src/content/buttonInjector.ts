// Closed shadow DOM injection so the page can't style or read the button.
// Each call returns the inner element you can attach a click handler to.

const BTN_LABEL = "Try it";
const ATTR = "data-tryit-injected";

export function injectButton(target: Element, onClick: () => void): void {
  if (target.hasAttribute(ATTR)) return;
  target.setAttribute(ATTR, "1");

  const host = document.createElement("span");
  host.style.marginLeft = "8px";
  host.style.verticalAlign = "middle";
  host.style.display = "inline-block";
  const root = host.attachShadow({ mode: "closed" });

  const style = document.createElement("style");
  style.textContent = `
    button {
      font: 600 11px/1 system-ui, sans-serif;
      padding: 4px 8px;
      border-radius: 6px;
      border: 1px solid #2a78c0;
      background: #2a78c0;
      color: white;
      cursor: pointer;
    }
    button:hover { background: #1f5e98; }
    button:focus { outline: 2px solid #88c0ff; outline-offset: 2px; }
    .badge {
      margin-left: 6px;
      font-size: 10px;
      opacity: 0.85;
    }
  `;
  const btn = document.createElement("button");
  btn.type = "button";
  btn.textContent = BTN_LABEL;
  const badge = document.createElement("span");
  badge.className = "badge";
  badge.textContent = "tryit";
  btn.appendChild(badge);
  btn.addEventListener("click", (ev) => {
    ev.preventDefault();
    ev.stopPropagation();
    onClick();
  });

  root.appendChild(style);
  root.appendChild(btn);
  target.appendChild(host);
}
