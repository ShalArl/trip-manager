#!/usr/bin/env python3
"""
Performance Report Generator

Analyzes Locust test results and generates a detailed performance report.
"""

import csv
import sys
from pathlib import Path
from datetime import datetime


def parse_stats_file(csv_file):
    """Parse Locust stats CSV file"""
    stats = {}

    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            name = row['Name']
            stats[name] = {
                'type': row['Type'],
                'num_requests': int(row['# requests']),
                'num_failures': int(row['# failures']),
                'median_response_time': int(row['Median response time']),
                'average_response_time': int(float(row['Average response time'])),
                'min_response_time': int(float(row['Min response time'])),
                'max_response_time': int(float(row['Max response time'])),
                'average_content_length': int(float(row['Average Content-Length'])),
                'requests_per_sec': float(row['Requests/sec']),
            }

    return stats


def generate_report(stats_file, output_file=None):
    """Generate a detailed performance report"""

    if not Path(stats_file).exists():
        print(f"❌ Stats file not found: {stats_file}")
        sys.exit(1)

    stats = parse_stats_file(stats_file)
    timestamp = datetime.now().isoformat()

    # Prepare output
    report = []
    report.append("=" * 80)
    report.append("🔍 Load Test Performance Report")
    report.append("=" * 80)
    report.append(f"Generated: {timestamp}")
    report.append(f"Source: {stats_file}")
    report.append("")

    # Overall metrics
    total_requests = sum(s['num_requests'] for s in stats.values())
    total_failures = sum(s['num_failures'] for s in stats.values())
    failure_rate = (total_failures / total_requests * 100) if total_requests > 0 else 0

    report.append("📊 Overall Metrics")
    report.append("-" * 80)
    report.append(f"  Total Requests:      {total_requests:,}")
    report.append(f"  Total Failures:      {total_failures:,}")
    report.append(f"  Failure Rate:        {failure_rate:.2f}%")
    report.append("")

    # Response time analysis
    avg_response_time = sum(s['average_response_time'] for s in stats.values()) / len(stats) if stats else 0
    max_response_time = max((s['max_response_time'] for s in stats.values()), default=0)
    min_response_time = min((s['min_response_time'] for s in stats.values()), default=0)

    report.append("⏱️  Response Times (ms)")
    report.append("-" * 80)
    report.append(f"  Min:                 {min_response_time}")
    report.append(f"  Max:                 {max_response_time}")
    report.append(f"  Average:             {avg_response_time:.0f}")
    report.append("")

    # Throughput
    total_rps = sum(s['requests_per_sec'] for s in stats.values())
    report.append(f"📈 Throughput")
    report.append("-" * 80)
    report.append(f"  Total RPS:           {total_rps:.2f}")
    report.append("")

    # Endpoint breakdown
    report.append("🔗 Endpoint Breakdown")
    report.append("-" * 80)
    report.append(f"{'Endpoint':<40} {'Requests':<12} {'Failures':<12} {'Avg Time':<12}")
    report.append("-" * 80)

    for name, stat in sorted(stats.items()):
        if name != 'Aggregated':
            endpoint = name[:39]
            report.append(
                f"{endpoint:<40} {stat['num_requests']:<12} "
                f"{stat['num_failures']:<12} {stat['average_response_time']:<12}"
            )

    report.append("")
    report.append("=" * 80)

    # Print report
    report_text = "\n".join(report)
    print(report_text)

    # Save to file if specified
    if output_file:
        with open(output_file, 'w') as f:
            f.write(report_text)
        print(f"\n✅ Report saved to: {output_file}")

    # Return metrics for programmatic use
    return {
        'total_requests': total_requests,
        'total_failures': total_failures,
        'failure_rate': failure_rate,
        'avg_response_time': avg_response_time,
        'max_response_time': max_response_time,
        'total_rps': total_rps,
    }


def check_thresholds(stats_file, response_time_threshold=1000, failure_rate_threshold=5):
    """Check if performance meets thresholds"""

    metrics = generate_report(stats_file)

    print("\n" + "=" * 80)
    print("🎯 Performance Thresholds Check")
    print("=" * 80)

    passed = True

    # Check response time
    if metrics['avg_response_time'] > response_time_threshold:
        print(f"❌ Response time {metrics['avg_response_time']:.0f}ms exceeds threshold {response_time_threshold}ms")
        passed = False
    else:
        print(f"✅ Response time {metrics['avg_response_time']:.0f}ms within threshold")

    # Check failure rate
    if metrics['failure_rate'] > failure_rate_threshold:
        print(f"❌ Failure rate {metrics['failure_rate']:.2f}% exceeds threshold {failure_rate_threshold}%")
        passed = False
    else:
        print(f"✅ Failure rate {metrics['failure_rate']:.2f}% within threshold")

    print("=" * 80)

    return 0 if passed else 1


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python analyze.py <stats_file> [output_file]")
        print("       python analyze.py --check <stats_file> [--response-time MS] [--failure-rate PERCENT]")
        sys.exit(1)

    if sys.argv[1] == '--check':
        if len(sys.argv) < 3:
            print("Usage: python analyze.py --check <stats_file> [--response-time MS] [--failure-rate PERCENT]")
            sys.exit(1)

        stats_file = sys.argv[2]
        response_threshold = 1000
        failure_threshold = 5

        # Parse optional thresholds
        for i in range(3, len(sys.argv), 2):
            if sys.argv[i] == '--response-time' and i + 1 < len(sys.argv):
                response_threshold = int(sys.argv[i + 1])
            elif sys.argv[i] == '--failure-rate' and i + 1 < len(sys.argv):
                failure_threshold = float(sys.argv[i + 1])

        exit_code = check_thresholds(stats_file, response_threshold, failure_threshold)
        sys.exit(exit_code)
    else:
        stats_file = sys.argv[1]
        output_file = sys.argv[2] if len(sys.argv) > 2 else None
        generate_report(stats_file, output_file)

