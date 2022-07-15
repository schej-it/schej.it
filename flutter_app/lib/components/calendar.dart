import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:linked_scroll_controller/linked_scroll_controller.dart';
import 'package:intl/intl.dart';
import 'package:intl/date_symbol_data_local.dart';

// Notes:
// For calendar view, make the times column its own list view,
// Make each day its own list view, and make the 3 day section a page view of list views
// Overlaying the events look at this video: https://www.youtube.com/watch?v=OOEyJ0ct0Sg

class Calendar extends StatefulWidget {
  const Calendar({Key? key}) : super(key: key);

  @override
  State<Calendar> createState() => _CalendarState();
}

class _CalendarState extends State<Calendar> {
  // Constants
  final double _timeColWidth = 50;
  final double _timeRowHeight = 45;
  final double _daySectionHeight = 62;

  // Controllers
  late final PageController _pageController;
  late final LinkedScrollControllerGroup _controllers;
  late final ScrollController _timeScrollController;

  // Other variables
  final DateTime _curDate = DateTime.now();
  final List<String> _timeStrings = <String>[];
  // Note: this _startDateOffset is hardcoded for now, i.e. if the user happens
  // to scroll back farther than 20 days, then they won't be able to scroll back
  // any farther
  final int _startDateOffset = -20;

  @override
  void initState() {
    super.initState();
    initializeDateFormatting();

    // Set up scroll controllers
    _controllers = LinkedScrollControllerGroup();
    _timeScrollController = _controllers.addAndGet();

    // Set initial scroll
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _controllers.jumpTo(8.25 * _timeRowHeight);
    });

    _pageController = PageController(
      viewportFraction: 1 / 3,
      initialPage: _startDateOffset.abs() + 1,
    );

    // Create a list of all the visible times, 1am - 11pm
    for (int i = 1; i < 24; ++i) {
      String timeText;
      if (i < 12) {
        timeText = '$i AM';
      } else if (i == 12) {
        timeText = '12 PM';
      } else {
        timeText = '${i - 12} PM';
      }
      _timeStrings.add(timeText);
    }
  }

  @override
  void dispose() {
    // Dispose all scroll controllers
    _pageController.dispose();
    _timeScrollController.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        _buildTimeColumn(),
        Expanded(child: _buildDaySection()),
      ],
    );
  }

  // Builds the section containing all the days in a horizontally scrolling
  // page view
  Widget _buildDaySection() {
    return FractionallySizedBox(
      heightFactor: 1,
      child: PageView.builder(
        controller: _pageController,
        itemBuilder: (BuildContext context, int index) {
          final date = _curDate.add(Duration(days: _startDateOffset + index));
          return _buildDay(date);
        },
      ),
    );
  }

  // Builds a column containing the given day and a scrollable list with
  // dividers representing the hour increments
  Widget _buildDay(DateTime date) {
    String dayText = DateFormat.E().format(date);
    int dateNum = date.day;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        SizedBox(
          height: _daySectionHeight,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(dayText,
                  style: date == _curDate
                      ? SchejFonts.body.copyWith(color: SchejColors.darkGreen)
                      : SchejFonts.body),
              Container(
                padding: const EdgeInsets.all(7),
                decoration: date == _curDate
                    ? const BoxDecoration(
                        color: SchejColors.darkGreen,
                        shape: BoxShape.circle,
                      )
                    : null,
                child: Text(dateNum.toString(),
                    style: date == _curDate
                        ? SchejFonts.header.copyWith(color: SchejColors.white)
                        : SchejFonts.header),
              ),
            ],
          ),
        ),
        const Divider(
          height: 1.15,
          thickness: 1.15,
          color: SchejColors.darkGray,
        ),
        Expanded(
          child: CalendarDay(
            controllers: _controllers,
            date: date,
            numRows: _timeStrings.length,
            rowHeight: _timeRowHeight,
          ),
        ),
      ],
    );
  }

  // Builds the column displaying all the times (1am-11pm) in a scroll view
  Widget _buildTimeColumn() {
    // itemBuilder for a row containing a time (e.g. 10am)
    Widget itemBuilder(BuildContext context, int index) {
      return SizedBox(
        height: _timeRowHeight,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.end,
          children: [
            // Time text
            Padding(
              padding: const EdgeInsets.only(right: 8.0),
              child: Align(
                alignment: Alignment.centerRight,
                child: Text(
                  _timeStrings[index],
                  style: SchejFonts.body.copyWith(color: SchejColors.darkGray),
                  textAlign: TextAlign.right,
                ),
              ),
            ),
            // Divider Fragment
            const SizedBox(
              width: 5,
              child: Divider(
                height: 1.15,
                color: SchejColors.lightGray,
                thickness: 1.15,
              ),
            ),
          ],
        ),
      );
    }

    return SizedBox(
      width: _timeColWidth,
      child: Column(
        children: [
          SizedBox(height: _daySectionHeight),
          const Divider(
            height: 1.15,
            thickness: 1.15,
            color: SchejColors.darkGray,
          ),
          Expanded(
            child: ListView.builder(
              scrollDirection: Axis.vertical,
              itemCount: _timeStrings.length,
              controller: _timeScrollController,
              physics: const BouncingScrollPhysics(
                  parent: AlwaysScrollableScrollPhysics()),
              itemBuilder: itemBuilder,
            ),
          ),
        ],
      ),
    );
  }
}

// Widget containing a list view with all the time dividers and events for the
// given day
class CalendarDay extends StatefulWidget {
  final LinkedScrollControllerGroup controllers;
  final DateTime date;
  final int numRows;
  final double rowHeight;

  const CalendarDay({
    Key? key,
    required this.controllers,
    required this.date,
    required this.numRows,
    required this.rowHeight,
  }) : super(key: key);

  @override
  State<CalendarDay> createState() => _CalendarDayState();
}

class _CalendarDayState extends State<CalendarDay> {
  late final ScrollController _controller;

  @override
  void initState() {
    super.initState();

    _controller = widget.controllers.addAndGet();
  }

  @override
  void dispose() {
    _controller.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      scrollDirection: Axis.vertical,
      itemCount: widget.numRows,
      controller: _controller,
      physics:
          const BouncingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
      itemBuilder: (BuildContext context, int index) {
        return SizedBox(
          height: widget.rowHeight,
          child: Row(
            children: const [
              VerticalDivider(
                width: 1.15,
                thickness: 1.15,
                color: SchejColors.lightGray,
              ),
              Expanded(
                child: Divider(
                  height: 1.15,
                  thickness: 1.15,
                  color: SchejColors.lightGray,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
