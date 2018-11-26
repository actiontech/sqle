using System.Threading.Tasks;
using Grpc.Core;
using SqlserverProto;

namespace SqlserverProtoServer
{
    public partial class SqlServerServiceImpl: SqlserverService.SqlserverServiceBase
    {
        public override Task<AdviseOutput> Advise(AdviseInput request, ServerCallContext context)
        {
            return base.Advise(request, context);
        }
    }
}
