import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/fonts.dart';

import '../router/app_router.gr.dart';

// Home page contains the bottom tab navigation bar as well as a scaffold to
// display the rest of the screens
class HomePage extends StatelessWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AutoTabsScaffold(
      routes: const [
        MySchejTabRoute(),
        // EventsTabRoute(),
        FriendsTabRoute(),
        ProfileTabRoute(),
      ],
      bottomNavigationBuilder: _buildBottomNavigation,
    );
  }

  Widget _buildBottomNavigation(_, tabsRouter) {
    return BottomNavigationBar(
      items: <BottomNavigationBarItem>[
        BottomNavigationBarItem(
          icon: tabsRouter.activeIndex == 0
              ? const Icon(Icons.calendar_today)
              : const Icon(Icons.calendar_today_outlined),
          label: 'My schej',
        ),
        // BottomNavigationBarItem(
        //   icon: tabsRouter.activeIndex == 1
        //       ? const Icon(MdiIcons.calendarBlankMultiple)
        //       : const Icon(Icomoon.calendarBlankMultipleOutline),
        //   label: 'My events',
        // ),
        BottomNavigationBarItem(
          icon: tabsRouter.activeIndex == 1
              ? const Icon(Icons.people)
              : const Icon(Icons.people_outline),
          label: 'Friends',
        ),
        BottomNavigationBarItem(
          icon: tabsRouter.activeIndex == 2
              ? const Icon(Icons.person)
              : const Icon(Icons.person_outline),
          label: 'Profile',
        ),
      ],
      type: BottomNavigationBarType.fixed,
      currentIndex: tabsRouter.activeIndex,
      onTap: tabsRouter.setActiveIndex,
      backgroundColor: SchejColors.white,
      selectedItemColor: SchejColors.darkGreen,
      unselectedItemColor: SchejColors.darkGray,
      selectedLabelStyle: SchejFonts.smallMedium.copyWith(color: null),
      unselectedLabelStyle: SchejFonts.smallMedium.copyWith(color: null),
    );
  }
}
