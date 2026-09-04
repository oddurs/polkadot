# Resolved on first use rather than at every shell start: `op read` blocks on
# 1Password auth, which stalls the prompt in every terminal.
function supabase_token --description 'Fetch and cache the Supabase access token from 1Password'
    if not set -q SUPABASE_ACCESS_TOKEN
        set -gx SUPABASE_ACCESS_TOKEN (op read 'op://Development/Supabase/access-token')
    end
    echo $SUPABASE_ACCESS_TOKEN
end
