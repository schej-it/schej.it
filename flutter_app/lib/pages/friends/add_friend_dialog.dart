import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/friends/add_friend_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:provider/provider.dart';

class AddFriendDialog extends StatefulWidget {
  final VoidCallback onClose;
  const AddFriendDialog({Key? key, required this.onClose}) : super(key: key);

  @override
  State<AddFriendDialog> createState() => _AddFriendDialogState();
}

class _AddFriendDialogState extends State<AddFriendDialog> {
  // Controllers
  final TextEditingController _searchTextController = TextEditingController();

  // Variables
  var results = [
    {'name': 'Winston Tilton', 'sent': false},
    {'name': 'Samantha Jones', 'sent': true},
    {'name': 'Tyler Smithson', 'sent': false},
    {'name': 'Arthi Singh', 'sent': false},
  ];

  var realResults = [];

  @override
  void initState() {
    super.initState();

    // Start listening to changes.
    _searchTextController.addListener(_updateSearchResults);
  }

  @override
  void dispose() {
    _searchTextController.dispose();
    super.dispose();
  }

  void _updateSearchResults() {
    ApiService api = context.read<ApiService>();
    api.refreshUserSearchResults(_searchTextController.text);
    print('Second text field: ${_searchTextController.text}');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(
        isRoot: true,
        underline: false,
        centerTitle: true,
        titleString: 'Add friends',
        actions: [
          IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            _buildSearchTextField(),
            Expanded(
              child: Consumer<ApiService>(
                builder: (context, api, child) => _buildResults(api),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 20),
      child: TextField(
        autofocus: true,
        controller: _searchTextController,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Search for a friend',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return textField;
  }

  Widget _buildResults(ApiService api) {
    return ListView.builder(
      itemCount: api.userSearchResults.length,
      itemBuilder: (context, index) {
        final result = api.userSearchResults[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: AddFriendCard(
            name: result.fullName,
            requestAlreadySent: false,
          ),
        );
      },
    );
  }
}
