#!/usr/bin/env python3
# requires python3.7+

import sys

from dataclasses import dataclass, field
from typing import List, Set


def SetData(data):
    # orders app
    orders_app = FakeApp("orders.log")
    orders_app.AddLogLine(
        2,
        "2020-06-02 14:33:56.531063 ORDERS price: 123.123, quantity: 1000, securityId: 999, side: BUY, bid: 124.0, ask: 125.0",
    )
    orders_app.AddLogLine(
        2,
        "2020-06-02 15:33:56.831063 ORDERS price: 140.123, quantity: 600, securityId: APPL, side: SELL, bid: 120.0, ask: 150.0",
    )
    data.AddApp(orders_app)

    # execution app
    execution_app = FakeApp("executions.log")
    execution_app.AddLogLine(
        3,
        "2020-06-03 18:53:54.531063 EXECUTION filledPrice: 532.21, filledQuantity: 5329, securityId: 1325, side: SELL, orderBid: 512.9, orderAsk: 532.8",
    )
    data.AddApp(execution_app)

    # position app
    position_app = FakeApp("positions.log")
    position_app.AddLogLine(
        3,
        "2020-07-12 01:54:23.124127 POSITION securityId: 154, netPosition: -1230, startOfDayPosition: 1240, dayTradedVolume: 2470",
    )
    data.AddApp(position_app)


@dataclass
class LogLine:
    Frequency: int
    Line: str


@dataclass
class FakeApp:
    FileName: str
    SortLines: bool = False
    LogLines: List[LogLine] = field(default_factory=list)

    def AddLogLine(self, frequency, line):
        self.LogLines.append(LogLine(frequency, line))


@dataclass
class FakeData:
    FakeApps: List[FakeApp] = field(default_factory=list)
    apps: Set[str] = field(init=False, default_factory=set)

    def AddApp(self, app):
        if app.FileName in self.apps:
            print(f"ERROR, App for file '{app.FileName}' already exist")
            sys.exit(1)
        self.apps.add(app.FileName)
        self.FakeApps.append(app)

    def WriteToFiles(self):
        for app in self.FakeApps:
            if len(app.LogLines) == 0:
                print(f"ERROR, No loglines provided for App '{app.FileName}'")
                continue
            if app.SortLines:
                app.LogLines.sort(key=lambda x: x.Line)
            with open(app.FileName, "w") as fo:
                for line in app.LogLines:
                    for i in range(line.Frequency):
                        print(line.Line, file=fo)
            print(f"{app.FileName} written")


def main():
    data = FakeData()
    SetData(data)
    data.WriteToFiles()


if __name__ == "__main__":
    main()
