using Grpc.Core;
using SqlserverProto;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.DependencyInjection;

namespace SqlserverProtoServer
{
    public class Program
    {
        // todo get Port from Main's args
        const int Port = 10086;

        public static async Task Main(string[] args)
        {
            var hostBuilder = new HostBuilder().ConfigureServices((hostContext, services) =>
            {
                Server server = new Server
                {
                    Services = { SqlserverService.BindService(new SqlServerServiceImpl()) },
                    Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
                };
                services.AddSingleton<Server>(server);
                services.AddSingleton<IHostedService, GrpcHostedService>();
            });
            await hostBuilder.RunConsoleAsync();
        }
    }

    public class GrpcHostedService : IHostedService
    {
        private Server _server;

        public GrpcHostedService(Server server)
        {
            _server = server;
        }

        public Task StartAsync(CancellationToken calcellationToken)
        {
            _server.Start();
            return Task.CompletedTask;
        }

        public async Task StopAsync(CancellationToken cancellation) => await _server.ShutdownAsync();
    }
}
