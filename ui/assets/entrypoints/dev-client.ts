import debug from "debug";

const log = debug("dev-client");
debug.enable("dev-client");

let reconnecting = false;

type ChangeEvent = {
  added: string[];
  removed: string[];
  updated: string[];
};

const urlMetaEl = document.head.querySelector<HTMLMetaElement>(
  '[name="dev-socket-url"]'
);
if (urlMetaEl === null) {
  throw new Error(
    'Expected a meta tag with name="dev-socket-url", but found none.'
  );
}
let socketURL: URL;
try {
  socketURL = new URL(urlMetaEl.content);
} catch (err) {
  throw new Error(`invalid socket URL: "${urlMetaEl.content}"`);
}

function connect() {
  log(
    `${reconnecting ? "reconnecting" : "connecting"} to dev server @ "${
      socketURL.href
    }"`
  );
  const socket = new WebSocket(socketURL, "ws");

  socket.onopen = function handleOpen(this: WebSocket, _event: Event) {
    if (reconnecting) {
      log("reconnected, reloading");
      location.reload();
    } else {
      log("connected");
    }
  };

  socket.onmessage = function handleMessage(
    this: WebSocket,
    event: MessageEvent<string>
  ) {
    const { added, removed, updated } = JSON.parse(event.data) as ChangeEvent;

    for (const link of document.getElementsByTagName("link")) {
      if (link.rel !== "stylesheet") continue;

      let url: URL | null = null;
      try {
        url = new URL(link.href);
      } catch (err) {
        console.error(`invalid stylesheet link href URL: "${link.href}"`);
        continue;
      }

      if (url.hostname === location.hostname) {
        for (const path of [...added, ...removed]) {
          if (path === url.pathname) {
            log('updated "%s"', path);
            location.reload();
            return;
          }
        }
      }
    }

    for (const script of document.getElementsByTagName("script")) {
      if (!script.src) continue;

      let url: URL | null = null;
      try {
        url = new URL(script.src);
      } catch (err) {
        console.error(`invalid script src URL: "${script.src}"`);
        continue;
      }

      if (url.hostname === location.hostname) {
        for (const path of [...added, ...removed, ...updated]) {
          if (path === url.pathname) {
            log('updated "%s"', path);
            location.reload();
            return;
          }
        }
      }
    }

    const links = Array.from(document.getElementsByTagName("link"));
    for (const link of links) {
      if (link.rel !== "stylesheet") continue;

      let url: URL | null = null;
      try {
        url = new URL(link.href);
      } catch (err) {
        console.error(`invalid stylesheet link href URL: "${link.href}"`);
      }
      if (url !== null && url.hostname === location.hostname) {
        const matching = updated.find((f) => url.pathname === f);
        if (matching) {
          log('updated "%s"', matching);
          const next = link.cloneNode() as HTMLLinkElement;
          url.searchParams.set("v", Math.random().toString(36).slice(2))
          next.href = url.href
          next.onload = () => {
            link.remove();
          };
          if (link.parentNode === null) {
            console.error("expected stylesheet link to have a parent node");
          } else {
            link.parentNode.insertBefore(next, link.nextSibling);
          }
        }
      }
    }
  };

  socket.onerror = function handleError(this: WebSocket, _event: Event) {
    // console.error("connection error", event);
  };

  socket.onclose = function handleClose(this: WebSocket, _event: CloseEvent) {
    // log("connection closed");
    reconnecting = true;
    setTimeout(connect, 2000);
  };
}

connect();
