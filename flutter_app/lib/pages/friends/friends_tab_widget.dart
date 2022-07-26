import 'package:flutter/material.dart';
import 'package:flutter_app/components/friend_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';

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
    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: SchejConstants.pagePadding.left,
      ),
      child: Column(
        children: [
          _buildSearchTextField(),
          Expanded(child: _buildFriendCards()),
        ],
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextField(
        controller: _searchTextController,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Search for a friend',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle.copyWith(color: SchejColors.black),
      ),
    );

    return textField;
  }

  Widget _buildFriendCards() {
    return ListView.builder(
      itemCount: friends.length,
      itemBuilder: (context, index) {
        final friend = friends[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: FriendCard(
            name: friend['name'] as String,
            status: friend['status'] as FriendStatus,
            showOverflowMenu: () {},
            curEventName: (friend['curEventName'] ?? '') as String,
          ),
        );
      },
    );
  }
}
