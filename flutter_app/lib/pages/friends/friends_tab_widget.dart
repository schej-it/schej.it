import 'package:flutter/material.dart';
import 'package:flutter_app/components/friend_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';

class FriendsTabWidget extends StatefulWidget {
  const FriendsTabWidget({Key? key}) : super(key: key);

  @override
  State<FriendsTabWidget> createState() => _FriendsTabWidgetState();
}

class _FriendsTabWidgetState extends State<FriendsTabWidget> {
  // Controllers
  late final TextEditingController _searchTextController;

  // Variables
  var friends = [
    {'name': 'Winston Tilton', 'status': FriendStatus.free},
    {
      'name': 'Winston Tilton',
      'status': FriendStatus.busy,
      'curEventName': 'BTG meeting'
    },
    {
      'name': 'Winston Tilton',
      'status': FriendStatus.busy,
      'curEventName': 'PSYC 336'
    },
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
  ];

  @override
  void initState() {
    super.initState();
    _searchTextController = TextEditingController();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildSearchTextField(),
        _buildFriendCards(),
      ],
    );
  }

  Widget _buildSearchTextField() {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextField(
        controller: _searchTextController,
        autocorrect: false,
        decoration: InputDecoration(
          border: OutlineInputBorder(
            borderRadius: SchejConstants.borderRadius,
          ),
          labelText: 'Search for a friend',
          prefixIcon: const Icon(Icons.search),
          prefixIconColor: SchejColors.black,
          focusColor: SchejColors.green,
          fillColor: Colors.red,
        ),
      ),
    );
  }

  Widget _buildFriendCards() {
    final friendCards = <Widget>[];
    for (final friend in friends) {
      friendCards.add(Padding(
        padding: const EdgeInsets.only(bottom: 8),
        child: FriendCard(
          name: friend['name'] as String,
          status: friend['status'] as FriendStatus,
          showOverflowMenu: () {},
          curEventName: (friend['curEventName'] ?? '') as String,
        ),
      ));
    }
    return Column(children: friendCards);
  }
}
