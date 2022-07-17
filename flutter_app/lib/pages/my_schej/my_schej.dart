import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar.dart';
import 'package:flutter_app/components/expand_transition.dart';
import 'package:flutter_app/components/month_calendar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/utils.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class MySchejPage extends StatefulWidget {
  const MySchejPage({Key? key}) : super(key: key);

  @override
  State<MySchejPage> createState() => _MySchejPageState();
}

class _MySchejPageState extends State<MySchejPage> {
  bool _monthSelector = false;
  bool _viewMenu = false;
  bool _eventNamesVisible = false;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);

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
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
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
                    calendarEvents: testCalendarEvents,
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
          const Text('July'),
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
        IconButton(
          icon: const Icon(MdiIcons.calendarBlankOutline),
          splashRadius: 15,
          onPressed: () => setState(() {
            _viewMenu = !_viewMenu;
          }),
        ),
        IconButton(
          icon: const Icon(MdiIcons.eye),
          splashRadius: 15,
          onPressed: () => setState(() {
            _eventNamesVisible = !_eventNamesVisible;
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
