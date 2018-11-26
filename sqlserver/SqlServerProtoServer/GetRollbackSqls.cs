using System.Threading.Tasks;
using Grpc.Core;
using SqlserverProto;

namespace SqlserverProtoServer 
{
    public partial class SqlServerServiceImpl: SqlserverService.SqlserverServiceBase 
    {
        public override Task<GetRollbackSqlsOutput> GetRollbackSqls(GetRollbackSqlsInput request, ServerCallContext context)
        {
            return base.GetRollbackSqls(request, context);
        }

    }
}