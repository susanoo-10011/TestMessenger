using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace AuthenticationService.Entity
{
    internal class TokenRecord
    {
        public long id { get; set; }
        public long user_id { get; set; }
        public string token { get; set; }
        public string ExpiresAt { get; set; }
        public string CreatedAt { get; set; }
    }
}
