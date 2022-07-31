import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar.dart';
import 'package:flutter_app/components/calendar_view_selector.dart';
import 'package:flutter_app/components/expand_transition.dart';
import 'package:flutter_app/components/month_calendar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/utils.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:path_provider/path_provider.dart';
import 'package:share_plus/share_plus.dart';

class MySchejPage extends StatefulWidget {
  const MySchejPage({Key? key}) : super(key: key);

  @override
  State<MySchejPage> createState() => _MySchejPageState();
}

class _MySchejPageState extends State<MySchejPage> {
  bool _monthSelector = false;
  int _daysVisible = 3;
  bool _eventTitlesVisible = true;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);
  final GlobalKey<CalendarState> _calendar = GlobalKey();

  void takeScreenshot() async {
    if (_calendar.currentState == null) return;

    final image = await _calendar.currentState!.getScreenshot();
    if (image == null) return;

    final directory = await getApplicationDocumentsDirectory();
    final imagePath = await File('${directory.path}/screenshot.png').create();
    await imagePath.writeAsBytes(image);
    await Share.shareFiles([imagePath.path]);
  }

  @override
  Widget build(BuildContext context) {
    final testCalendarEvents = CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Event',
          startDate: getDateWithTime(DateTime.now(), 9.5),
          endDate: getDateWithTime(DateTime.now(), 12),
        ),
        CalendarEvent(
          title: 'Introduction to Failure Analysis',
          startDate: getDateWithTime(DateTime.now(), 13),
          endDate: getDateWithTime(DateTime.now(), 14.5),
        ),
        CalendarEvent(
          title: 'War',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 15),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 20),
        ),
      ],
    );

    return Scaffold(
      appBar: _buildAppBar(),
      floatingActionButton: FloatingActionButton(
        onPressed: takeScreenshot,
        backgroundColor: SchejColors.darkGreen,
        child: const Icon(MdiIcons.share),
      ),
      body: Container(
        color: SchejColors.white,
        // padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            ExpandTransition(
              visible: _monthSelector,
              child: Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: MonthCalendar(
                  selectedDay: _selectedDay,
                  onDaySelected: (selectedDay) => setState(() {
                    _selectedDay = selectedDay;
                  }),
                ),
              ),
            ),
            Expanded(
              child: Stack(
                children: [
                  Calendar(
                    key: _calendar,
                    calendarEvents: testCalendarEvents,
                    daysVisible: _daysVisible,
                    eventTitlesVisible: _eventTitlesVisible,
                    selectedDay: _selectedDay,
                    onDaySelected: (selectedDay) => setState(() {
                      _selectedDay = selectedDay;
                    }),
                  ),
                  if (_monthSelector) _buildGestureDetector(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return SchejAppBar(
      title: Row(
        children: [
          Text(
            DateFormat.MMMM().format(_selectedDay),
            style: SchejFonts.header,
          ),
          IconButton(
            icon: _monthSelector
                ? const Icon(MdiIcons.chevronUp)
                : const Icon(MdiIcons.chevronDown),
            splashRadius: 15,
            onPressed: () => setState(() {
              _monthSelector = !_monthSelector;
            }),
          ),
        ],
      ),
      actions: [
        CalendarViewSelector(
          onSelected: (int value) => setState(() {
            _daysVisible = value;
          }),
        ),
        IconButton(
          icon: _eventTitlesVisible
              ? const Icon(MdiIcons.eye)
              : const Icon(MdiIcons.eyeOff),
          splashRadius: 15,
          onPressed: () => setState(() {
            _eventTitlesVisible = !_eventTitlesVisible;
          }),
        ),
      ],
      underline: false,
      isRoot: true,
    );
  }

  // Builds a gesture detector enabling the user to swipe up to dismiss the month
  // calendar
  double _totalDrag = 0;
  Widget _buildGestureDetector() {
    return GestureDetector(
      onVerticalDragUpdate: (details) {
        if (details.primaryDelta != null) {
          _totalDrag += details.primaryDelta!;
        }

        if (_totalDrag < -20) {
          setState(() {
            _totalDrag = 0;
            _monthSelector = false;
          });
        }
      },
      child: Container(
        color: Colors.transparent,
      ),
    );
  }
}
