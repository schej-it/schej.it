import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/utils.dart';
import 'package:linked_scroll_controller/linked_scroll_controller.dart';
import 'package:intl/intl.dart';
import 'package:screenshot/screenshot.dart';

// TODO: fix bug where if you pinch to zoom, the time scroll controller gets
// out of sync with the individual day scroll controllers

// The Calendar widget contains a widget to view the user's daily events
class Calendar extends StatefulWidget {
  final CalendarEvents calendarEvents;
  final DateTime selectedDay;
  final void Function(DateTime) onDaySelected;
  final int daysVisible;
  final bool eventTitlesVisible;

  const Calendar({
    Key? key,
    required this.calendarEvents,
    required this.selectedDay,
    required this.onDaySelected,
    this.daysVisible = 3,
    this.eventTitlesVisible = true,
  }) : super(key: key);

  @override
  State<Calendar> createState() => CalendarState();
}

class CalendarState extends State<Calendar> {
  // Constants
  final double _timeColWidth = 60;
  final double _daySectionHeight = 62;
  final double _minTimeRowHeight = 25;
  final double _maxTimeRowHeight = 90;

  // Controllers
  late PageController _pageController;
  late final LinkedScrollControllerGroup _controllers;
  late final ScrollController _timeScrollController;
  final ScreenshotController _screenshotController = ScreenshotController();

  // Other variables
  final DateTime _curDate = getDateWithTime(DateTime.now(), 0);
  final List<String> _timeStrings = <String>[];
  // Note: this _startDateOffset is hardcoded for now, i.e. if the user happens
  // to scroll back farther than 365 days, then they won't be able to scroll back
  // any farther
  final int _startDateOffset = -365;
  bool _pageControllerAnimating = false;
  bool _isTakingScreenshot = false;
  double _timeRowHeight = 45;

  @override
  void initState() {
    super.initState();

    // Set up scroll controllers
    _controllers = LinkedScrollControllerGroup();
    _timeScrollController = _controllers.addAndGet();

    // Set initial scroll
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _controllers.jumpTo(8.25 * _timeRowHeight);
    });

    // Set up page controller
    _pageController = PageController(
      viewportFraction: 1 / widget.daysVisible,
      initialPage: _startDateOffset.abs() + 1,
    );
    _pageController.addListener(_pageControllerListener);

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
  void didUpdateWidget(covariant Calendar oldWidget) {
    super.didUpdateWidget(oldWidget);

    // If selectedDay changed, animate the page to the correct day
    if (widget.selectedDay != oldWidget.selectedDay) {
      _animateToDay(widget.selectedDay, oldDay: oldWidget.selectedDay);
    }

    // Change the number of days visible
    if (widget.daysVisible != oldWidget.daysVisible) {
      _pageController =
          PageController(viewportFraction: 1 / widget.daysVisible);
      _pageController.addListener(_pageControllerListener);
    }
  }

  // Returns a Uint8List containing a screenshot of the user's schej
  Future<Uint8List?> getScreenshot() async {
    final oldPage = _pageController.page!.truncate();
    setState(() {
      _isTakingScreenshot = true;
    });
    final image = await _screenshotController.capture();
    setState(() {
      _isTakingScreenshot = false;
    });

    // HACK: reassign pagecontroller so it doesn't reset the page it's on
    _pageController = PageController(
      initialPage: oldPage,
      viewportFraction: 1 / widget.daysVisible,
    );

    return image;
  }

  // Animate the page to the given day
  void _animateToDay(DateTime day, {DateTime? oldDay}) async {
    Duration diff = day.difference(_curDate);
    int index = diff.inDays - _startDateOffset + 1;

    // Scale animation duration with difference between day and oldDay
    Duration animDuration = const Duration(milliseconds: 1000);
    if (oldDay != null) {
      int dayDiff = day.difference(oldDay).inDays.abs();
      animDuration = Duration(milliseconds: 270 + 30 * dayDiff);
    }

    _pageControllerAnimating = true;
    await _pageController.animateToPage(
      index,
      duration: animDuration,
      curve: Curves.easeInOut,
    );
    _pageControllerAnimating = false;
  }

  // Listener for whenever the page controller changes
  void _pageControllerListener() {
    // Update selectedDay whenever the page changes
    if (!_pageControllerAnimating &&
        _pageController.page != null &&
        _pageController.page!.truncate() == _pageController.page) {
      int newOffset = _pageController.page!.truncate() + _startDateOffset - 1;
      DateTime newDay = _curDate.add(Duration(days: newOffset));
      if (newDay != widget.selectedDay) {
        widget.onDaySelected(newDay);
      }
    }
  }

  double _oldTimeRowHeight = 0;
  void _onScaleStart(ScaleStartDetails details) {
    _oldTimeRowHeight = _timeRowHeight;
  }

  void _onScaleUpdate(ScaleUpdateDetails details) {
    setState(() {
      _timeRowHeight = _oldTimeRowHeight * details.verticalScale;
      if (_timeRowHeight < _minTimeRowHeight) {
        _timeRowHeight = _minTimeRowHeight;
      } else if (_timeRowHeight > _maxTimeRowHeight) {
        _timeRowHeight = _maxTimeRowHeight;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isTakingScreenshot) {
      return _buildScreenshot();
    }

    return GestureDetector(
      onScaleStart: _onScaleStart,
      onScaleUpdate: _onScaleUpdate,
      child: Row(
        children: [
          _buildTimeColumn(_timeScrollController),
          Expanded(child: _buildDaySection(_pageController)),
        ],
      ),
    );
  }

  // Builds the section containing all the days in a horizontally scrolling
  // page view
  Widget _buildDaySection(PageController controller) {
    return FractionallySizedBox(
      heightFactor: 1,
      child: PageView.builder(
        controller: controller,
        itemBuilder: (BuildContext context, int index) {
          int dayOffset =
              _startDateOffset + index - 1 + (widget.daysVisible - 1) ~/ 2;
          DateTime utcDate = _curDate.add(Duration(days: dayOffset));
          // Need to convert to local date in order to get all the events for
          // the local day
          DateTime localDate = getLocalDayFromUtcDay(utcDate);
          return _buildDay(localDate);
        },
      ),
    );
  }

  // Builds a column containing the given day and a scrollable list with
  // dividers representing the hour increments
  Widget _buildDay(DateTime date) {
    String dayText = DateFormat.E().format(date);
    int dateNum = date.day;
    bool isCurDate = date == getLocalDayFromUtcDay(_curDate);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        SizedBox(
          height: _daySectionHeight,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(dayText,
                  style: isCurDate
                      ? SchejFonts.body.copyWith(color: SchejColors.darkGreen)
                      : SchejFonts.body),
              Container(
                padding: const EdgeInsets.all(7),
                decoration: isCurDate
                    ? const BoxDecoration(
                        color: SchejColors.darkGreen,
                        shape: BoxShape.circle,
                      )
                    : null,
                child: Text(dateNum.toString(),
                    style: isCurDate
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
            key: ObjectKey(date),
            controllers: _controllers,
            date: date,
            events: widget.calendarEvents.eventsByDay[date],
            eventTitlesVisible: widget.eventTitlesVisible,
            numRows: _timeStrings.length,
            rowHeight: _timeRowHeight,
          ),
        ),
      ],
    );
  }

  // Builds the column displaying all the times (1am-11pm) in a scroll view
  Widget _buildTimeColumn(ScrollController controller) {
    // itemBuilder for a row containing a time (e.g. 10am)
    Widget itemBuilder(BuildContext context, int index) {
      // Account for the half hour before and half hour after for 12am
      if (index == 0 || index == _timeStrings.length + 1) {
        return SizedBox(height: _timeRowHeight / 2);
      }

      return SizedBox(
        key: ValueKey(_timeStrings[index - 1]),
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
                  _timeStrings[index - 1],
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
              key: const ValueKey('timeColumn'),
              scrollDirection: Axis.vertical,
              itemCount: _timeStrings.length + 2,
              controller: controller,
              physics: const BouncingScrollPhysics(
                  parent: AlwaysScrollableScrollPhysics()),
              itemBuilder: itemBuilder,
            ),
          ),
        ],
      ),
    );
  }

  // Builds the screenshot view when taking a screenshot
  Widget _buildScreenshot() {
    final localTimeScrollController = ScrollController(
      initialScrollOffset: _timeScrollController.offset,
    );
    final localPageController = PageController(
      initialPage: _pageController.page!.truncate(),
      viewportFraction: 1 / widget.daysVisible,
    );

    return SingleChildScrollView(
      controller: ScrollController(initialScrollOffset: 73),
      child: Screenshot(
        controller: _screenshotController,
        child: Container(
          height: MediaQuery.of(context).size.height,
          color: SchejColors.white,
          child: Column(
            children: [
              _buildScreenshotHeader(),
              Expanded(
                child: Row(
                  children: [
                    _buildTimeColumn(localTimeScrollController),
                    Expanded(child: _buildDaySection(localPageController)),
                  ],
                ),
              ),
              _buildScreenshotFooter(),
            ],
          ),
        ),
      ),
    );
  }

  // Builds the screenshot header containing the name of the user + date ranges
  Widget _buildScreenshotHeader() {
    final startDateString = DateFormat.MMMMd().format(widget.selectedDay);
    final endDateString = DateFormat.MMMMd()
        .format(widget.selectedDay.add(Duration(days: widget.daysVisible - 1)));
    final dateString = widget.daysVisible == 1
        ? startDateString
        : '$startDateString - $endDateString';

    return Container(
      padding: const EdgeInsets.only(top: 8, left: 8, right: 8, bottom: 20),
      child: Column(
        children: [
          const Text('{USERNAME}\'s schej', style: SchejFonts.header),
          Text(dateString, style: SchejFonts.body),
        ],
      ),
    );
  }

  // Builds the screenshot footer containing timezone + logo
  Widget _buildScreenshotFooter() {
    return Container(
      padding: const EdgeInsets.all(8),
      alignment: Alignment.centerLeft,
      child: Text('Timezone: ${DateTime.now().timeZoneName}'),
    );
  }
}

// Widget containing a list view with all the time dividers and events for the
// given day
class CalendarDay extends StatefulWidget {
  final LinkedScrollControllerGroup controllers;
  final DateTime date;
  final List<CalendarEvent>? events;
  final int numRows;
  final double rowHeight;
  final bool eventTitlesVisible;

  const CalendarDay({
    Key? key,
    required this.controllers,
    required this.date,
    required this.events,
    required this.numRows,
    required this.rowHeight,
    this.eventTitlesVisible = true,
  }) : super(key: key);

  @override
  State<CalendarDay> createState() => _CalendarDayState();
}

class _CalendarDayState extends State<CalendarDay> {
  // Controllers
  late final ScrollController _emptyController;
  late final ScrollController _timeRowsController;

  // Variables
  final LayerLink _layerLink = LayerLink();

  @override
  void initState() {
    super.initState();

    _emptyController = widget.controllers.addAndGet();
    _timeRowsController = widget.controllers.addAndGet();
  }

  @override
  void dispose() {
    _emptyController.dispose();
    _timeRowsController.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          _buildEmpty(),
          _buildEvents(),
          _buildTimeRows(),
        ],
      ),
    );
  }

  // Builds a list view containing the events for this day
  Widget _buildEvents() {
    return Stack(
      children: widget.events == null
          ? []
          : widget.events!
              .map((event) => CalendarEventWidget(
                    event: event,
                    hourHeight: widget.rowHeight,
                    layerLink: _layerLink,
                    titleVisible: widget.eventTitlesVisible,
                  ))
              .toList(),
    );
  }

  // Builds the listview containing the dividers representing the time rows
  Widget _buildTimeRows() {
    final timeRows = ListView.builder(
      scrollDirection: Axis.vertical,
      itemCount: widget.numRows + 2,
      controller: _timeRowsController,
      physics:
          const BouncingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
      itemBuilder: (BuildContext context, int index) {
        // Account for the half hour before and half hour after for 12am
        if (index == 0 || index == widget.numRows + 1) {
          return Align(
            alignment: Alignment.centerLeft,
            child: SizedBox(
              height: widget.rowHeight / 2,
              child: const VerticalDivider(
                width: 1.15,
                thickness: 1.15,
                color: SchejColors.lightGray,
              ),
            ),
          );
        }

        final divider = SizedBox(
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

        return divider;
      },
    );

    return timeRows;
  }

  // Builds an empty scrollable widget that takes up the entire screen for the
  // event to use as an offset
  Widget _buildEmpty() {
    return SingleChildScrollView(
      scrollDirection: Axis.vertical,
      controller: _emptyController,
      physics:
          const BouncingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
      child: CompositedTransformTarget(
        link: _layerLink,
        child: SizedBox(height: widget.numRows * widget.rowHeight),
      ),
    );
  }
}

// CalendarEventWidget is a graphical representation of a user's calendar event
// The layerLink allows us to change the position of the event according to the
// current scroll value
class CalendarEventWidget extends StatefulWidget {
  final CalendarEvent event;
  final double hourHeight;
  final LayerLink layerLink;
  final bool titleVisible;

  const CalendarEventWidget({
    Key? key,
    required this.event,
    required this.hourHeight,
    required this.layerLink,
    this.titleVisible = true,
  }) : super(key: key);

  @override
  State<CalendarEventWidget> createState() => _CalendarEventWidgetState();
}

class _CalendarEventWidgetState extends State<CalendarEventWidget> {
  @override
  Widget build(BuildContext context) {
    // It looks like the stack is clipping the calendar event widgets for some reason
    return FractionallySizedBox(
      widthFactor: 1,
      child: CompositedTransformFollower(
        link: widget.layerLink,
        showWhenUnlinked: false,
        offset: Offset(0, widget.event.startTime * widget.hourHeight),
        child: Container(
          margin: const EdgeInsets.only(right: 2),
          padding: const EdgeInsets.only(top: 7, right: 7, left: 7),
          height: (widget.event.endTime - widget.event.startTime) *
              widget.hourHeight,
          decoration: BoxDecoration(
            color: widget.titleVisible
                ? SchejColors.lightGreen
                : SchejColors.fadedGreen,
            borderRadius: const BorderRadius.all(Radius.circular(5)),
          ),
          child: widget.titleVisible
              ? Text(
                  widget.event.title,
                  style: SchejFonts.body.copyWith(color: SchejColors.white),
                )
              : Text(
                  'BUSY',
                  style: SchejFonts.body.copyWith(color: SchejColors.darkGreen),
                ),
        ),
      ),
    );
  }
}
