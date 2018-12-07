using Grpc.Core;
using SqlserverProto;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.DependencyInjection;
using CommandLine;
using IniParser;
using IniParser.Model;
using System.Diagnostics;
using System;

namespace SqlserverProtoServer {
    public class Options {
        [Option('p', "port", Required = false, HelpText = "grpc server port")]
        public int Port { get; set; }
        [Option('i', "pidfile", Required = false, HelpText = "pid file")]
        public string Pidfile { get; set; }
        [Option('c', "config", Required = true, HelpText = "config file")]
        public string Config { get; set; }
    }

    public class Program {
        // default grpc port
        public static int Port = 10001;
        // defaulr pidfile
        public static string PidFile = "sqled_sqlserver.pid";

        public static async Task Main(string[] args) {
            Parser.Default.ParseArguments<Options>(args)
                  .WithParsed<Options>(o => {
                      if (o.Port > 0) {
                          Port = o.Port;

                      }
                      if (o.Pidfile != "") {
                          PidFile = o.Pidfile;
                      }

                      if (o.Config != "") {
                          var parser = new FileIniDataParser();
                          IniData iniData = parser.ReadFile(o.Config);
                          string portStr = iniData["server"]["port"];
                          if (portStr != "") {
                              Port = Int32.Parse(portStr);
                          }
                      }

                  });

            var hostBuilder = new HostBuilder().ConfigureServices((hostContext, services) => {
                Server server = new Server {
                    Services = { SqlserverService.BindService(new SqlServerServiceImpl()) },
                    Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
                };

                services.AddSingleton<Server>(server);
                services.AddSingleton<IHostedService, GrpcHostedService>();
            });
            await hostBuilder.RunConsoleAsync();
        }
    }

    public class GrpcHostedService : IHostedService {
        private Server _server;

        public GrpcHostedService(Server server) {
            _server = server;
        }

        public Task StartAsync(CancellationToken calcellationToken) {
            int processID = Process.GetCurrentProcess().Id;
            if (System.IO.File.Exists(Program.PidFile)) {
                throw new Exception(String.Format("There has a pidfile:{0}", System.IO.File.ReadAllText(Program.PidFile)));
            }
            System.IO.File.WriteAllText(Program.PidFile, String.Format("{0}", processID));

            _server.Start();
            return Task.CompletedTask;
        }

        public async Task StopAsync(CancellationToken cancellation) {
            System.IO.File.Delete(Program.PidFile);
            await _server.ShutdownAsync();
        }
    }
}
