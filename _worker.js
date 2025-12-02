//Server
import { connect } from 'cloudflare:sockets';

export default {
  /**
   * @param {Request} req
   * @param {Object} env
   * @param {ExecutionContext} ctx
   */
  async fetch(req, env, ctx) {
    try {
      // 1. 验证 WebSocket 请求
      if (req.headers.get('Upgrade') !== 'websocket') {
        return new Response('Not WebSocket', { status: 426 });
      }

      // 2. 验证密码 (优先使用环境变量 PASSWORD，否则使用默认值)
      const secret = env.PASSWORD || 'testPASSword';
      if (req.headers.get('X-Password') !== secret) {
        return new Response('Unauthorized', { status: 403 });
      }

      // 3. 验证目标地址格式 (简化的逻辑验证替代冗长正则)
      const target = req.headers.get('X-Target') || '';
      const [hostname, portStr] = target.split(':');
      const port = parseInt(portStr);
      if (!hostname || isNaN(port) || port < 1 || port > 65535) {
        return new Response('Invalid Target', { status: 400 });
      }

      // 4. 建立连接
      const [client, server] = Object.values(new WebSocketPair());
      const socket = connect({ hostname, port });

      server.accept();

      // 5. 定义流传输逻辑 (高效管道)
      // TCP -> WebSocket
      const tcpToWs = socket.readable.pipeTo(new WritableStream({
        write(chunk) { server.send(chunk); },
        close() { server.close(); }
      }));

      // WebSocket -> TCP
      const wsToTcp = new ReadableStream({
        start(ctrl) {
          server.addEventListener('message', ({ data }) => {
            if (data instanceof ArrayBuffer) ctrl.enqueue(new Uint8Array(data));
          });
          server.addEventListener('close', () => ctrl.close());
          server.addEventListener('error', e => ctrl.error(e));
        }
      }).pipeTo(socket.writable);

      // 6. 保持连接直到结束
      ctx.waitUntil(
        Promise.all([tcpToWs, wsToTcp])
          .catch(() => {})
          .finally(() => server.close())
      );

      return new Response(null, { status: 101, webSocket: client });

    } catch (err) {
      return new Response(err.message, { status: 500 });
    }
  }
};
