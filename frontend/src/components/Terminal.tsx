import { AttachAddon } from "@xterm/addon-attach";
import { FitAddon } from "@xterm/addon-fit";
import { server } from "@go/models";
import { useEffect, useState } from "react";
import { useXTerm } from "react-xtermjs";

export default function Terminal(): JSX.Element {
  const { instance, ref } = useXTerm();
  const fitAddon = new FitAddon();
  const [tty, setTty] = useState<string | null>(null);

  useEffect(() => {
    if (tty === null) {
      return;
    }
    const websocket = new WebSocket(`ws://localhost:8080/${tty}`);
    const attachAddon = new AttachAddon(websocket);
    instance?.loadAddon(fitAddon);
    instance?.loadAddon(attachAddon);
    instance?.onResize(async (c) => {
      const req = new server.SetTtySizeDto({ row: c.rows, col: c.cols, tty });
      fetch("http://localhost:8080/resize", {
        method: "POST",
        body: JSON.stringify(req),
      });
    });
    const handleResize = () => fitAddon.fit();
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [fitAddon, instance, ref, tty]);

  useEffect(() => {
    const getTty = async () => {
      const resp = server.GetTtyDto.createFrom(
        await (await fetch("http://localhost:8080/tty")).json(),
      );
      setTty(resp.tty);
    };
    getTty();
  }, []);

  useEffect(() => {
    if (fitAddon) {
      setTimeout(() => fitAddon.fit(), 0);
    }
  }, [fitAddon]);

  return <div style={{ height: "100%", width: "100%" }} ref={ref}></div>;
}
