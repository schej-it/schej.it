import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      title: 'schej',
      home: NavBarView(),
    );
  }
}

class NavBarView extends StatefulWidget {
  const NavBarView({Key? key}) : super(key: key);

  @override
  State<NavBarView> createState() => _NavBarViewState();
}

class _NavBarViewState extends State<NavBarView> {
  int _selectedIndex = 0;

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('test')),
      bottomNavigationBar: BottomNavigationBar(
        items: <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: _selectedIndex == 0
                ? const Icon(Icons.calendar_today)
                : const Icon(Icons.calendar_today_outlined),
            label: 'My schej',
          ),
          BottomNavigationBarItem(
            icon: _selectedIndex == 1
                ? const Icon(MdiIcons.calendarMultipleCheck)
                : const Icon(MdiIcons.calendarBlankMultiple),
            label: 'My events',
          ),
          BottomNavigationBarItem(
            icon: _selectedIndex == 2
                ? const Icon(Icons.people)
                : const Icon(Icons.people_outline),
            label: 'Friends',
          ),
          BottomNavigationBarItem(
            icon: _selectedIndex == 3
                ? const Icon(Icons.person)
                : const Icon(Icons.person_outline),
            label: 'Profile',
          ),
        ],
        type: BottomNavigationBarType.fixed,
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
        backgroundColor: Colors.white,
        selectedItemColor: Colors.black,
        unselectedItemColor: Colors.grey,
      ),
    );
  }
}
