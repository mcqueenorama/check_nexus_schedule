Check the number of nexus scheduled tasks 
=========================================

    Usage of ./check_nexus_schedule:
      -c=20: critical level for numner of scheduled tasks
      -h="http://nuch.com/nexus": base url for nexus  like http://zup.com/nexus
      -v=false: verbose output
      -w=10: warning level for number of scheduled tasks

    Creds comes from ~/.netrc.

Build it:
---------

  go build

or:

  go build check_nexus_schedule.go


