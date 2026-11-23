#!/usr/bin/env python3
"""
Summarizes the RESULT lines written by bench_run.sh into Markdown tables
(mean and standard deviation of the throughput for every configuration).

Usage: summarize.py <modules results file> [<pool results file>]
"""
import collections
import re
import statistics
import sys

RESULT_RE = re.compile(
    r"RESULT (?P<name>\S+) logging=(?P<logging>\d) rules=(?P<rules>\d) ml=(?P<ml>\d) "
    r"processes=(?P<processes>\d+) requests=\d+ time=\S+ rps=(?P<rps>[\d.]+) status=(?P<status>.*)"
)


def read_results(path):
    results = collections.defaultdict(list)
    invalid = []
    with open(path) as f:
        for line in f:
            match = RESULT_RE.match(line.strip())
            if not match:
                continue
            if match["status"] != "OK":
                invalid.append(line.strip())
                continue
            key = (match["logging"], match["rules"], match["ml"], int(match["processes"]))
            results[key].append(float(match["rps"]))
    return results, invalid


def mean_std(values):
    std = statistics.stdev(values) if len(values) > 1 else 0.0
    return f"{statistics.mean(values):,.0f} ± {std:,.0f}"


def check(flag):
    return "✓" if flag == "1" else ""


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(1)

    all_invalid = []

    modules, invalid = read_results(sys.argv[1])
    all_invalid += invalid
    if modules:
        print("## Throughput by active module configuration\n")
        print("| Proxy | Logging | Rule engine | ML classifier | Runs | Req/s |")
        print("|---|---|---|---|---|---|")
        for (logging, rules, ml, _), values in modules.items():
            print(f"| ✓ | {check(logging)} | {check(rules)} | {check(ml)} | {len(values)} | {mean_std(values)} |")

    if len(sys.argv) > 2:
        pool, invalid = read_results(sys.argv[2])
        all_invalid += invalid
        print("\n## Throughput by number of classifier processes\n" if modules else "## Throughput by number of classifier processes\n")
        print("| Processes | ML classifier only (req/s) | All modules (req/s) |")
        print("|---|---|---|")
        for processes in sorted({key[3] for key in pool}):
            ml_only = pool.get(("0", "0", "1", processes))
            all_modules = pool.get(("1", "1", "1", processes))
            print(f"| {processes} | {mean_std(ml_only) if ml_only else '-'} | {mean_std(all_modules) if all_modules else '-'} |")

    if all_invalid:
        print("\n## Invalid runs (excluded from the tables)\n")
        for line in all_invalid:
            print(f"- `{line}`")
        sys.exit(2)


if __name__ == "__main__":
    main()
