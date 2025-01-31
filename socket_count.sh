#!/bin/bash

declare -A TCP_STATES=(
    ["00"]="UNKNOWN"
    ["01"]="ESTABLISHED"
    ["02"]="SYN_SENT"
    ["03"]="SYN_RECV"
    ["04"]="FIN_WAIT1"
    ["05"]="FIN_WAIT2"
    ["06"]="TIME_WAIT"
    ["07"]="CLOSE"
    ["08"]="CLOSE_WAIT"
    ["09"]="LAST_ACK"
    ["0A"]="LISTEN"
    ["0B"]="CLOSING"
    ["0C"]="NEW_SYN_RECV"
)

declare -A state_count
for key in "${TCP_STATES[@]}"; do
    state_count["$key"]=0
done

parse_tcp_table() {
    local file="$1"

    while read -r line; do
        state=$(echo "$line" | awk '{print $4}')
        [[ -z "$state" ]] && continue

        readable_state=${TCP_STATES[$state]:-"UNKNOWN"}

        ((state_count["$readable_state"]++))
    done < <(tail -n +2 "$file")
}

if [[ ! -f /proc/net/tcp || ! -f /proc/net/tcp6 ]]; then
    echo "Error: /proc/net/tcp or /proc/net/tcp6 not found!"
    exit 1
fi

parse_tcp_table /proc/net/tcp
parse_tcp_table /proc/net/tcp6

echo "=== TCP Connection States Count ==="
for state in "${!state_count[@]}"; do
    echo "$state: ${state_count[$state]}"
done
