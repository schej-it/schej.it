import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar/calendar.dart';
import 'package:flutter_app/components/calendar/calendar_view_selector.dart';
import 'package:flutter_app/components/expand_transition.dart';
import 'package:flutter_app/components/month_calendar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/utils.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:path_provider/path_provider.dart';
import 'package:provider/provider.dart';
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
  bool _isTakingScreenshot = false;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);
  final GlobalKey<CalendarState> _calendar = GlobalKey();

  void takeScreenshot() async {
    setState(() => _isTakingScreenshot = true);

    if (_calendar.currentState == null) return;

    final image = await _calendar.currentState!.getScreenshot();
    if (image == null) return;

    final directory = await getApplicationDocumentsDirectory();
    final imagePath = await File('${directory.path}/screenshot.png').create();
    await imagePath.writeAsBytes(image);
    await Share.shareFiles([imagePath.path]);

    setState(() => _isTakingScreenshot = false);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      floatingActionButton: Padding(
        padding: const EdgeInsets.only(bottom: 12),
        child: FloatingActionButton(
          heroTag: 'schejFab',
          onPressed: _isTakingScreenshot ? null : takeScreenshot,
          backgroundColor: SchejColors.darkGreen,
          child: const Icon(MdiIcons.share, size: 28),
        ),
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
            Consumer<ApiService>(
              builder: (context, api, child) => Expanded(
                child: Stack(
                  children: [
                    Calendar(
                      mode: CalendarMode.schej,
                      key: _calendar,
                      userIds: {api.authUser!.id},
                      daysVisible: _daysVisible,
                      showEventTitles: _eventTitlesVisible,
                      selectedDay: _selectedDay,
                      onDaySelected: (selectedDay) => setState(() {
                        _selectedDay = selectedDay;
                      }),
                    ),
                    if (_monthSelector) _buildGestureDetector(),
                  ],
                ),
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
            DateFormat.yMMMM().format(_selectedDay),
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
        Padding(
          padding: const EdgeInsets.only(right: 15),
          child: IconButton(
            icon: _eventTitlesVisible
                ? const Icon(MdiIcons.eye)
                : const Icon(MdiIcons.eyeOff),
            splashRadius: 15,
            onPressed: () => setState(() {
              _eventTitlesVisible = !_eventTitlesVisible;
            }),
          ),
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
