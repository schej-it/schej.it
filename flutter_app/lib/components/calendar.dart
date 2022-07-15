import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:linked_scroll_controller/linked_scroll_controller.dart';
import 'package:intl/intl.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:sorted_list/sorted_list.dart';

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

  // Controllers
  late final PageController _pageController;
  late final LinkedScrollControllerGroup _controllers;
  late final ScrollController _timeScrollController;
  late final Map<DateTime, ScrollController> _dayScrollControllers =
      <DateTime, ScrollController>{};

  // Other variables
  final DateTime _curDate = DateTime.now();
  final SortedList<DateTime> _loadedDates =
      SortedList<DateTime>((a, b) => a.compareTo(b));
  final List<String> _timeStrings = <String>[];

  @override
  void initState() {
    super.initState();
    initializeDateFormatting();

    // Set up scroll controllers
    _controllers = LinkedScrollControllerGroup();
    _timeScrollController = _controllers.addAndGet();
    for (int i = -3; i <= 3; ++i) {
      final date = _curDate.add(Duration(days: i));
      _loadedDates.add(date);
      _dayScrollControllers[date] = _controllers.addAndGet();
    }
    // _controllers.jumpTo(8.25 * _timeRowHeight);

    _pageController = PageController(
      viewportFraction: 1 / 3,
      initialPage: _loadedDates.indexOf(_curDate),
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
    for (var controller in _dayScrollControllers.values) {
      controller.dispose();
    }

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
    final days = <Widget>[];
    for (final date in _loadedDates) {
      days.add(_buildDay(date));
    }

    return FractionallySizedBox(
      heightFactor: 1,
      child: PageView(
        controller: _pageController,
        children: days,
      ),
    );
  }

  // Builds a column containing the given day and a scrollable list with
  // dividers representing the hour increments
  Widget _buildDay(DateTime date) {
    String dayText = DateFormat.E().format(date);
    int dateNum = date.day;

    final hourIncrements = ListView.builder(
      scrollDirection: Axis.vertical,
      itemCount: _timeStrings.length,
      controller: _dayScrollControllers[date],
      itemBuilder: (BuildContext context, int index) {
        return SizedBox(
          height: _timeRowHeight,
          child: Row(
            children: const [
              VerticalDivider(
                width: 1.15,
                thickness: 1.15,
                color: SchejColors.lightGray,
              ),
              Expanded(
                child: Divider(
                  thickness: 1.15,
                  color: SchejColors.lightGray,
                ),
              ),
            ],
          ),
        );
      },
    );

    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        SizedBox(
          height: 50,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(dayText, style: SchejFonts.body),
              Text(dateNum.toString(), style: SchejFonts.header),
            ],
          ),
        ),
        const Divider(color: SchejColors.darkGray),
        Expanded(child: hourIncrements),
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
          const SizedBox(height: 50),
          const Divider(color: SchejColors.darkGray),
          Expanded(
            child: ListView.builder(
              scrollDirection: Axis.vertical,
              itemCount: _timeStrings.length,
              controller: _timeScrollController,
              itemBuilder: itemBuilder,
            ),
          ),
        ],
      ),
    );
  }
}
