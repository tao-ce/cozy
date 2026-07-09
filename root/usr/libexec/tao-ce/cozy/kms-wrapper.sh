#!/usr/bin/bash
#

exec /usr/bin/kmscon \
    --login \
    --switchvt \
    --vt "$TTY" \
    -- /usr/bin/bash -ci \
        "$*; kill -15 \$PPID;"

