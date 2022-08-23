import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_app/components/user_avatar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/availabilities.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/utils.dart';
import 'package:linked_scroll_controller/linked_scroll_controller.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:screenshot/screenshot.dart';

// TODO: fix bug where if you pinch to zoom, the time scroll controller gets
// out of sync with the individual day scroll controllers

// The Calendar widget contains a widget to view the user's daily events
class Calendar extends StatefulWidget {
  final Set<String> userIds;
  final DateTime selectedDay;
  final void Function(DateTime) onDaySelected;
  final int daysVisible;
  final bool showEventTitles;
  final bool showAvatars;
  final bool showAvailability;
  final String? activeUserId;

  const Calendar({
    Key? key,
    required this.userIds,
    required this.selectedDay,
    required this.onDaySelected,
    this.daysVisible = 3,
    this.showEventTitles = true,
    this.showAvatars = false,
    this.showAvailability = false,
    this.activeUserId,
  }) : super(key: key);

  @override
  State<Calendar> createState() => CalendarState();
}

class CalendarState extends State<Calendar> {
  //
  // Constants
  //
  final double _timeColWidth = 60;
  final double _daySectionHeight = 62;
  final double _minTimeRowHeight = 25;
  final double _maxTimeRowHeight = 90;

  //
  // Controllers
  //
  late PageController _pageController;
  late final LinkedScrollControllerGroup _controllers;
  late final ScrollController _timeScrollController;
  final ScreenshotController _screenshotController = ScreenshotController();

  //
  // Other variables
  //
  final DateTime _curDate = getDateWithTime(DateTime.now(), 0);
  final List<String> _timeStrings = <String>[];
  late final Map<String, CalendarEvents> _calendarEvents =
      <String, CalendarEvents>{};
  final Map<String, DayRange> _loadedDayRanges = <String, DayRange>{};
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

    // Create calendarEvents map
    _updateCalendarEvents();
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

    _updateCalendarEvents();
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

  // Updates the calendarEvents map and dates loaded based on the current day
  void _updateCalendarEvents() {
    const int rangeRadius = 7;
    DateTime left =
        widget.selectedDay.subtract(const Duration(days: rangeRadius));
    DateTime right = widget.selectedDay.add(const Duration(days: rangeRadius));

    for (String userId in widget.userIds) {
      DayRange dayRangeToLoad;

      if (_loadedDayRanges[userId] == null) {
        // Load the initial set of calendar events
        _loadedDayRanges[userId] = DayRange(start: left, end: right);
        dayRangeToLoad = DayRange(start: left, end: right);
      } else {
        // Determine the direction in which to extend the day range
        DayRange loadedDayRange = _loadedDayRanges[userId]!;
        String extendDirection = '';
        if (!loadedDayRange.isInRange(left) &&
            !loadedDayRange.isInRange(right)) {
          if (loadedDayRange.compareToDay(left) < 0) {
            // Day range is to the left of local range
            extendDirection = 'right';
          } else {
            // Day range is to the right of local range
            extendDirection = 'left';
          }
        } else if (!loadedDayRange.isInRange(left)) {
          extendDirection = 'left';
        } else if (!loadedDayRange.isInRange(right)) {
          extendDirection = 'right';
        } else {
          continue;
        }

        // Extend the day range in that direction
        if (extendDirection == 'left') {
          final newStart = left.subtract(const Duration(days: rangeRadius));
          dayRangeToLoad = DayRange(
            start: loadedDayRange.start.subtract(const Duration(days: 1)),
            end: newStart,
          );
          _loadedDayRanges[userId]!.start = newStart;
        } else if (extendDirection == 'right') {
          final newEnd = right.add(const Duration(days: rangeRadius));
          dayRangeToLoad = DayRange(
            start: loadedDayRange.end.add(const Duration(days: 1)),
            end: newEnd,
          );
          _loadedDayRanges[userId]!.end = newEnd;
        } else {
          continue;
        }
      }

      // Load the specified day range
      _loadDayRangeForUser(dayRangeToLoad, userId);
    }
  }

  // Loads all calendar events for the given day range for the given user
  Future<void> _loadDayRangeForUser(DayRange dayRange, String userId) async {
    final api = context.read<ApiService>();
    final events = await api.getCalendarEvents(
      id: userId,
      startDate: dayRange.start,
      endDate: dayRange.end,
    );

    setState(() {
      if (_calendarEvents[userId] == null) {
        _calendarEvents[userId] = CalendarEvents(events: events);
      } else {
        _calendarEvents[userId]!.addEvents(events);
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

    Map<String, List<CalendarEvent>> events = _calendarEvents.map(
        (id, calendarEvents) =>
            MapEntry(id, calendarEvents.getEventsForDay(date)));

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
            events: events,
            showEventTitles: widget.showEventTitles,
            showAvatars: widget.showAvatars,
            showAvailability: widget.showAvailability,
            activeUserId: widget.activeUserId,
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
  final Map<String, List<CalendarEvent>> events;
  final int numRows;
  final double rowHeight;
  final bool showEventTitles;
  final bool showAvatars;
  final bool showAvailability;
  final String? activeUserId;

  const CalendarDay({
    Key? key,
    required this.controllers,
    required this.date,
    required this.events,
    required this.numRows,
    required this.rowHeight,
    this.showEventTitles = true,
    this.showAvatars = false,
    this.showAvailability = false,
    this.activeUserId,
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
    final children = <Widget>[];

    if (widget.showAvailability) {
      final availabilities =
          Availabilities.getUsersAvailabilityForDay(widget.date, widget.events);
      children.addAll(availabilities
          .map((a) => AvailabilityBlockWidget(
                availabilityBlock: a,
                maxNumUsersAvailable: widget.events.keys.length,
                hourHeight: widget.rowHeight,
                layerLink: _layerLink,
              ))
          .toList());

      // Only show current user's schedule if there's an active user
      if (widget.activeUserId != null) {
        children.addAll(widget.events[widget.activeUserId]!
            .map((event) => CalendarEventWidget(
                  event: event,
                  hourHeight: widget.rowHeight,
                  layerLink: _layerLink,
                  showTitle: widget.showEventTitles,
                  showAvatar: widget.showAvatars,
                  userId: widget.activeUserId,
                  activeUserId: widget.activeUserId,
                ))
            .toList());
      }
    } else {
      List<String> userIds = List.from(widget.events.keys);
      if (widget.activeUserId != null) {
        // Make sure activeUser is the last user to be displayed. This ensures
        // activeUser's events are at the very front
        userIds.remove(widget.activeUserId!);
        userIds.add(widget.activeUserId!);
      }

      for (String userId in userIds) {
        if (widget.events[userId] != null) {
          children.addAll(widget.events[userId]!
              .map((event) => CalendarEventWidget(
                    event: event,
                    hourHeight: widget.rowHeight,
                    layerLink: _layerLink,
                    showTitle: widget.showEventTitles,
                    showAvatar: widget.showAvatars,
                    userId: userId,
                    activeUserId: widget.activeUserId,
                    // marginLeftPercent: .20 * i,
                  ))
              .toList());
        }
      }
    }

    return Stack(
      children: children,
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

  final bool showTitle;
  final bool showBusy;
  final bool showAvatar;
  final bool noBorder;
  final String? userId;
  final String? activeUserId;
  final double marginLeftPercent;
  final Color? backgroundColor;

  const CalendarEventWidget({
    Key? key,
    required this.event,
    required this.hourHeight,
    required this.layerLink,
    this.showTitle = true,
    // Whether to show "BUSY" when title not shown
    this.showBusy = true,
    this.showAvatar = false,
    this.noBorder = false,
    this.userId,
    this.activeUserId,
    this.marginLeftPercent = 0,
    this.backgroundColor,
  }) : super(key: key);

  @override
  State<CalendarEventWidget> createState() => _CalendarEventWidgetState();
}

class _CalendarEventWidgetState extends State<CalendarEventWidget> {
  late final Color _containerColor;
  late final Color _textColor;

  @override
  void initState() {
    super.initState();

    // Determine containerColor and textColor
    if (widget.backgroundColor != null) {
      _containerColor = widget.backgroundColor!;
      _textColor = SchejColors.white;
    } else {
      if (widget.showTitle &&
          (widget.activeUserId == null ||
              widget.userId == widget.activeUserId)) {
        _containerColor = SchejColors.lightBlue;
        _textColor = SchejColors.white;
      } else {
        _containerColor = SchejColors.fadedLightBlue;
        _textColor = SchejColors.lightBlue;
      }
    }
  }

  String _getProfileImage(ApiService api) {
    if (widget.userId == null) return '';

    if (widget.userId! == api.authUser!.id) return api.authUser!.picture;

    User? friend = api.getFriendById(widget.userId!);
    if (friend == null) return '';

    return friend.picture;
  }

  // Build a CompositedTransformFollower that scrolls with the scrollview
  @override
  Widget build(BuildContext context) {
    return CompositedTransformFollower(
      link: widget.layerLink,
      showWhenUnlinked: false,
      offset: Offset(0, widget.event.startTime * widget.hourHeight),
      child: Align(
        alignment: Alignment.topRight,
        child: FractionallySizedBox(
          widthFactor: 1 - widget.marginLeftPercent,
          child: Padding(
            padding: const EdgeInsets.only(right: 2),
            child: Stack(
              clipBehavior: Clip.none,
              children: [
                _buildContainer(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // Build the main container for the event, with the event text and time block
  // TODO: this crashes when endTime is on the next day!!! (so if somebody has a late event)
  Widget _buildContainer() {
    double height =
        (widget.event.endTime - widget.event.startTime) * widget.hourHeight;
    if (height < 0) {
      height = 0;
    }

    return Positioned(
      top: 0,
      left: 0,
      right: 0,
      height: height,
      child: Stack(
        children: [
          Container(
            padding: const EdgeInsets.only(top: 7, right: 7, left: 7),
            height: height,
            width: double.infinity,
            decoration: BoxDecoration(
              color: _containerColor,
              borderRadius: const BorderRadius.all(Radius.circular(5)),
              border: !widget.noBorder
                  ? Border.all(
                      color: SchejColors.white,
                      width: 2,
                    )
                  : null,
            ),
            child: widget.showTitle
                ? Text(
                    widget.event.title,
                    style: SchejFonts.body.copyWith(color: _textColor),
                  )
                : Text(
                    widget.showBusy ? 'BUSY' : '',
                    style: SchejFonts.body.copyWith(color: _textColor),
                  ),
          ),
          if (widget.showAvatar) _buildAvatar(),
        ],
      ),
    );
  }

  // Build the avatar of the user that owns this event
  Widget _buildAvatar() {
    return Positioned(
      right: 0,
      bottom: 0,
      child: Padding(
        padding: const EdgeInsets.all(3),
        child: Consumer<ApiService>(
          builder: (context, api, child) => UserAvatar(
            src: _getProfileImage(api),
            radius: 10,
          ),
        ),
      ),
    );
  }
}

class AvailabilityBlockWidget extends StatefulWidget {
  final AvailabilityBlock availabilityBlock;
  final int maxNumUsersAvailable;
  final double hourHeight;
  final LayerLink layerLink;
  final double marginLeftPercent;

  const AvailabilityBlockWidget({
    Key? key,
    required this.availabilityBlock,
    required this.maxNumUsersAvailable,
    required this.hourHeight,
    required this.layerLink,
    this.marginLeftPercent = 0,
  }) : super(key: key);

  @override
  State<AvailabilityBlockWidget> createState() =>
      _AvailabilityBlockWidgetState();
}

class _AvailabilityBlockWidgetState extends State<AvailabilityBlockWidget> {
  @override
  Widget build(BuildContext context) {
    double alpha = (widget.availabilityBlock.usersAvailable.length /
            widget.maxNumUsersAvailable) *
        255;
    // Scale alpha if not everybody is available to make the times where everybody
    // is available stand out more
    if (alpha.round() != 255) {
      alpha *= 0.9;
    }
    final backgroundColor = SchejColors.green.withAlpha(alpha.round());

    return CalendarEventWidget(
      event: CalendarEvent(
        title: '',
        startDate: widget.availabilityBlock.startDate,
        endDate: widget.availabilityBlock.endDate,
      ),
      backgroundColor: backgroundColor,
      hourHeight: widget.hourHeight,
      layerLink: widget.layerLink,
      showTitle: false,
      showBusy: false,
      noBorder: true,
    );
  }
}
