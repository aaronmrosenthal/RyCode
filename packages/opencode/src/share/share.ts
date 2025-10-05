import { Bus } from "../bus"
import { Installation } from "../installation"
import { Session } from "../session"
import { MessageV2 } from "../session/message-v2"
import { Log } from "../util/log"

export namespace Share {
  const log = Log.create({ service: "share" })

  let queue: Promise<void> = Promise.resolve()
  const pending = new Map<string, any>()

  export async function sync(key: string, content: any) {
    const [root, ...splits] = key.split("/")
    if (root !== "session") return
    const [sub, sessionID] = splits
    if (sub === "share") return
    const share = await Session.getShare(sessionID).catch(() => {})
    if (!share) return
    const { secret } = share
    pending.set(key, content)
    queue = queue
      .then(async () => {
        const content = pending.get(key)
        if (content === undefined) return

        try {
          // BUG FIX: Add 30s timeout to prevent indefinite hangs
          const controller = new AbortController()
          const timeout = setTimeout(() => controller.abort(), 30000)

          try {
            const response = await fetch(`${URL}/share_sync`, {
              method: "POST",
              body: JSON.stringify({
                sessionID: sessionID,
                secret,
                key: key,
                content,
              }),
              signal: controller.signal,
            })

            log.info("synced", {
              key: key,
              status: response.status,
            })
          } finally {
            clearTimeout(timeout)
          }
        } catch (error: any) {
          // BUG FIX: Log errors instead of silently swallowing them
          log.error("sync failed", {
            key,
            error: error.message,
            type: error.name,
          })
        } finally {
          // BUG FIX: Always delete from pending to prevent memory leak
          pending.delete(key)
        }
      })
  }

  export function init() {
    Bus.subscribe(Session.Event.Updated, async (evt) => {
      await sync("session/info/" + evt.properties.info.id, evt.properties.info)
    })
    Bus.subscribe(MessageV2.Event.Updated, async (evt) => {
      await sync("session/message/" + evt.properties.info.sessionID + "/" + evt.properties.info.id, evt.properties.info)
    })
    Bus.subscribe(MessageV2.Event.PartUpdated, async (evt) => {
      await sync(
        "session/part/" +
          evt.properties.part.sessionID +
          "/" +
          evt.properties.part.messageID +
          "/" +
          evt.properties.part.id,
        evt.properties.part,
      )
    })
  }

  export const URL =
    process.env["OPENCODE_API"] ??
    (Installation.isSnapshot() || Installation.isDev() ? "https://api.dev.opencode.ai" : "https://api.opencode.ai")

  export async function create(sessionID: string) {
    // BUG FIX: Add 30s timeout to prevent indefinite hangs
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 30000)

    try {
      return await fetch(`${URL}/share_create`, {
        method: "POST",
        body: JSON.stringify({ sessionID: sessionID }),
        signal: controller.signal,
      })
        .then((x) => x.json())
        .then((x) => x as { url: string; secret: string })
    } finally {
      clearTimeout(timeout)
    }
  }

  export async function remove(sessionID: string, secret: string) {
    // BUG FIX: Add 30s timeout to prevent indefinite hangs
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 30000)

    try {
      return await fetch(`${URL}/share_delete`, {
        method: "POST",
        body: JSON.stringify({ sessionID, secret }),
        signal: controller.signal,
      }).then((x) => x.json())
    } finally {
      clearTimeout(timeout)
    }
  }
}
